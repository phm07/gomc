package player

import (
	"fmt"
	"github.com/google/uuid"
	"gomc/src/connection"
	"gomc/src/protocol"
	"gomc/src/protocol/packet"
	"gomc/src/textcomponent"
)

type Player struct {
	Conn       *connection.Connection
	UUID       uuid.UUID
	Name       string
	X, Y, Z    float64
	Yaw, Pitch float32
}

func (p *Player) SendMessage(msg *textcomponent.Component) error {
	if p.Conn.State != protocol.StatePlay {
		return fmt.Errorf("connection not in play state")
	}
	return p.Conn.SendPacket(&packet.ClientboundPlaySystemMessage{
		Content: msg.MarshalNBT().Marshal(),
		Overlay: false,
	})
}
