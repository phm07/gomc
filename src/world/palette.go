package world

import (
	"bytes"
	"gomc/src/protocol/types"
)

type Palette interface {
	Marshal() []byte
}

type PaletteSingleValued struct {
	Value types.VarInt
}

func (p *PaletteSingleValued) Marshal() []byte {
	return p.Value.Marshal()
}

type PaletteIndirect struct {
	Length  types.VarInt
	Palette []types.VarInt
}

func (p *PaletteIndirect) Marshal() []byte {
	var buf bytes.Buffer
	buf.Write(p.Length.Marshal())
	for _, v := range p.Palette {
		buf.Write(v.Marshal())
	}
	return buf.Bytes()
}

func (p *PaletteIndirect) IndexOf(v types.VarInt) int {
	for i, u := range p.Palette {
		if v == u {
			return i
		}
	}
	return -1
}

type PaletteDirect struct{}

func (*PaletteDirect) Marshal() []byte {
	return []byte{}
}
