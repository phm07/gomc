package protocol

import (
	"bytes"
	"gomc/src/protocol/types"
)

type Packet struct {
	Length   types.VarInt
	PacketID types.VarInt
	Data     *bytes.Buffer
}

func (p *Packet) Marshal() []byte {
	p.Length = types.VarInt(p.Data.Len() + p.PacketID.Len())
	var buf bytes.Buffer
	buf.Write(p.Length.Marshal())
	buf.Write(p.PacketID.Marshal())
	buf.Write(p.Data.Bytes())
	return buf.Bytes()
}
