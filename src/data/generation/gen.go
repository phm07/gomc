package main

import (
	"encoding/json"
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/samber/lo"
	"golang.org/x/exp/maps"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

//go:generate bash prepare.sh
//go:generate go run gen.go

type PropType int

const (
	Bool PropType = iota
	Int
	Enum
)

type State struct {
	Id         int               `json:"id"`
	Default    bool              `json:"default"`
	Properties map[string]string `json:"properties"`
}

type Block struct {
	Properties map[string][]string `json:"properties"`
	States     []State             `json:"states"`
}

type OrderedBlock struct {
	SmallestStateID int
	Name            string
	Properties      map[string][]string `json:"properties"`
	States          []State             `json:"states"`
}

func main() {
	bytes, err := os.ReadFile("generated/reports/blocks.json")
	if err != nil {
		panic(err)
	}

	var blocks map[string]Block
	if err := json.Unmarshal(bytes, &blocks); err != nil {
		panic(err)
	}

	var ordered []OrderedBlock
	for name, block := range blocks {
		smallestId := math.MaxInt
		for _, state := range block.States {
			smallestId = min(smallestId, state.Id)
		}
		ordered = append(ordered, OrderedBlock{
			SmallestStateID: smallestId,
			Name:            name,
			Properties:      block.Properties,
			States:          block.States,
		})
	}

	slices.SortFunc(ordered, func(a, b OrderedBlock) int {
		if a.SmallestStateID < b.SmallestStateID {
			return -1
		}
		if a.SmallestStateID > b.SmallestStateID {
			return 1
		}
		return 0
	})

	f := jen.NewFile("data")
	f.HeaderComment("Code generated by generation/gen.go. DO NOT EDIT.")

	for _, block := range ordered {
		namePas := snakeCaseToPascalCase(block.Name)
		propTypes := make(map[string]PropType)
		var props []jen.Code
		for prop, values := range block.Properties {
			propPas := snakeCaseToPascalCase(prop)
			if slices.Equal(values, []string{"true", "false"}) {
				props = append(props, jen.Id(propPas).Bool())
				propTypes[prop] = Bool
			} else if areNumbers(values) {
				props = append(props, jen.Comment(fmt.Sprintf("Valid values: %s", strings.Join(values, ", "))).
					Line().Id(propPas).Int())
				propTypes[prop] = Int
			} else {
				enumName := namePas + propPas
				f.Type().Id(enumName).String().Line()
				f.Const().DefsFunc(func(g *jen.Group) {
					for _, value := range values {
						g.Id(snakeCaseToPascalCase(enumName + snakeCaseToPascalCase(value))).Id(enumName).Op("=").Lit(value)
					}
				}).Line()
				props = append(props, jen.Id(propPas).Id(enumName))
				propTypes[prop] = Enum
			}
		}

		f.Type().Id(namePas).Struct(props...).Line()

		g := &jen.Group{}
		makeSwitchCase(g, namePas, lo.Keys(block.Properties), block.Properties, propTypes, make(map[string]string), block.States)

		if len(block.Properties) > 0 {
			def := 0
			for _, s := range block.States {
				if s.Default {
					def = s.Id
				}
			}
			g.Line().Return(jen.Lit(def)).Comment("default state")
		}

		f.Func().Params(jen.Id("x").Id(namePas)).Id("Id").Params().Uint16().Block(g).Line()
	}

	err = f.Save("../zz_block_states.go")
	if err != nil {
		panic(err)
	}
}

func makeSwitchCase(j *jen.Group, namePas string, toAdd []string, values map[string][]string, types map[string]PropType,
	valuesSoFar map[string]string, states []State) {
	if len(toAdd) == 0 {
		var s State
		for _, state := range states {
			match := true
			for k, v := range state.Properties {
				if valuesSoFar[k] != v {
					match = false
					break
				}
			}
			if match {
				s = state
				break
			}
		}
		j.Return(jen.Lit(s.Id))
		return
	}
	currProp := toAdd[0]
	currPropPas := snakeCaseToPascalCase(currProp)
	toAdd = toAdd[1:]
	var cases []jen.Code
	for _, value := range values[currProp] {
		m := maps.Clone(valuesSoFar)
		m[currProp] = value
		var c jen.Code
		switch types[currProp] {
		case Bool:
			c = jen.Lit(value == "true")
		case Int:
			v, _ := strconv.Atoi(value)
			c = jen.Lit(v)
		case Enum:
			id := namePas + currPropPas + snakeCaseToPascalCase(value)
			c = jen.Id(id)
		}
		cases = append(cases, jen.Case(c).BlockFunc(func(g *jen.Group) {
			makeSwitchCase(g, namePas, toAdd, values, types, m, states)
		}))
	}
	j.Switch(jen.Id("x").Dot(currPropPas)).Block(cases...)
}

func snakeCaseToPascalCase(s string) string {
	s = strings.TrimPrefix(s, "minecraft:")
	var sb strings.Builder
	iscap := true
	for _, r := range s {
		if r == '_' {
			iscap = true
			continue
		}
		if iscap {
			r = unicode.ToUpper(r)
		}
		sb.WriteRune(r)
		iscap = false
	}
	return sb.String()
}

func areNumbers(s []string) bool {
	for _, v := range s {
		if _, err := strconv.Atoi(v); err != nil {
			return false
		}
	}
	return true
}
