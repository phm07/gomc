package packet

import (
	"bytes"
	"errors"
	"gomc/src/protocol"
	"gomc/src/protocol/types"
)

// this packet is so complex that it needs custom serialization

type ClientboundPlayPlayerInfoUpdateAction interface {
	Serialize() []byte
	Mask() byte
}

type ClientboundPlayPlayerInfoUpdateActionAddPlayerProperty struct {
	Name      types.String
	Value     types.String
	IsSigned  types.Boolean
	Signature types.String
}

func (p *ClientboundPlayPlayerInfoUpdateActionAddPlayerProperty) Serialize() []byte {
	var buf bytes.Buffer
	buf.Write(p.Name.Marshal())
	buf.Write(p.Value.Marshal())
	buf.Write(p.IsSigned.Marshal())
	if p.IsSigned {
		buf.Write(p.Signature.Marshal())
	}
	return buf.Bytes()
}

type ClientboundPlayPlayerInfoUpdateActionAddPlayer struct {
	Name       types.String
	Properties []*ClientboundPlayPlayerInfoUpdateActionAddPlayerProperty
}

func (a *ClientboundPlayPlayerInfoUpdateActionAddPlayer) Serialize() []byte {
	var buf bytes.Buffer
	buf.Write(a.Name.Marshal())
	buf.Write(types.VarInt(len(a.Properties)).Marshal())
	for _, p := range a.Properties {
		buf.Write(p.Serialize())
	}
	return buf.Bytes()
}

func (*ClientboundPlayPlayerInfoUpdateActionAddPlayer) Mask() byte {
	return 0x01
}

type ClientboundPlayPlayerInfoUpdatePlayer struct {
	UUID    types.UUID
	Actions []ClientboundPlayPlayerInfoUpdateAction
}

type ClientboundPlayPlayerInfoUpdate struct {
	Players []*ClientboundPlayPlayerInfoUpdatePlayer
}

func (p *ClientboundPlayPlayerInfoUpdate) Serialize() []byte {
	var buf bytes.Buffer
	var mask byte
	if len(p.Players) == 0 {
		panic("cannot serialize ClientboundPlayPlayerInfoUpdate with no players")
	}
	for _, a := range p.Players[0].Actions {
		mask |= a.Mask()
	}
	buf.WriteByte(mask)
	buf.Write(types.VarInt(len(p.Players)).Marshal())
	for _, p := range p.Players {
		buf.Write(p.UUID.Marshal())
		for _, a := range p.Actions {
			buf.Write(a.Serialize())
		}
	}
	return buf.Bytes()
}

func (p *ClientboundPlayPlayerInfoUpdate) Deserialize(_ []byte) error {
	return errors.New("not implemented")
}

func (p *ClientboundPlayPlayerInfoUpdate) State() protocol.State {
	return protocol.StatePlay
}

func (p *ClientboundPlayPlayerInfoUpdate) ID() int {
	return 0x3c
}
