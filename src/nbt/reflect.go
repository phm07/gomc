package nbt

import (
	"reflect"
	"slices"
	"strings"
)

func Marshal(v any) Tag {
	val := reflect.ValueOf(v)
	t := reflect.TypeOf(v)
	switch val.Kind() {
	case reflect.Struct:
		tag := &CompoundTag{}
		for i := 0; i < t.NumField(); i++ {
			fVal := val.Field(i)
			ft := t.Field(i)
			name := ft.Tag.Get("nbt")
			if name == "" {
				name = ft.Tag.Get("json")
			}
			if name == "" {
				name = ft.Name
			}
			var omitempty bool
			if idx := strings.IndexByte(name, ','); idx != -1 {
				tags := strings.Split(name[idx+1:], ",")
				name = name[:idx]
				omitempty = slices.Contains(tags, "omitempty")
			}
			if omitempty && fVal.IsZero() {
				continue
			}
			elem := Marshal(fVal.Interface())
			// elem.Name = name
			reflect.Indirect(reflect.ValueOf(elem)).FieldByName("Name").SetString(name)
			tag.Data = append(tag.Data, elem)
		}
		return tag

	case reflect.Map:
		tag := &CompoundTag{}
		r := val.MapRange()
		for r.Next() {
			key, val := r.Key(), r.Value()
			elem := Marshal(val.Interface())
			// elem.Name = name
			reflect.Indirect(reflect.ValueOf(elem)).FieldByName("Name").SetString(key.String())
			tag.Data = append(tag.Data, elem)
		}
		return tag

	case reflect.Slice:
		tag := &ListTag[Tag]{}
		for i := 0; i < val.Len(); i++ {
			tag.Data = append(tag.Data, Marshal(val.Index(i).Interface()))
		}
		return tag

	case reflect.Int:
		i := val.Interface().(int)
		return &IntTag{Data: int32(i)}

	case reflect.Int64:
		i := val.Interface().(int64)
		return &LongTag{Data: i}

	case reflect.String:
		s := val.Interface().(string)
		return &StringTag{Data: s}

	case reflect.Float32:
		f := val.Interface().(float32)
		return &FloatTag{Data: f}

	case reflect.Float64:
		f := val.Interface().(float64)
		return &DoubleTag{Data: f}

	case reflect.Uint8:
		b := val.Interface().(uint8)
		return &ByteTag{Data: b}

	default:
		panic("cannot marshal " + val.Kind().String())
	}
}
