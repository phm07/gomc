package server

import (
	"encoding/json"
	"gomc/src/connection"
	"gomc/src/protocol/packet"
	"gomc/src/protocol/types"
	"gomc/src/status"
)

func (s *Server) handleStatus(c *connection.Connection, p packet.SerializablePacket) error {
	switch p := p.(type) {

	case *packet.ServerboundStatusRequest:
		stat := status.GetStatus()
		statBytes, err := json.Marshal(stat)
		if err != nil {
			return err
		}

		return c.SendPacket(&packet.ClientboundStatusResponse{
			Json: types.String(statBytes),
		})

	case *packet.ServerboundStatusPing:
		return c.SendPacket(&packet.ClientboundStatusPong{Payload: p.Payload})

	default:
		return nil
	}
}
