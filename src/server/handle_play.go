package server

import (
	"gomc/src/connection"
	"gomc/src/protocol/packet"
)

func (s *Server) handlePlay(c *connection.Connection, p packet.SerializablePacket) error {
	switch p := p.(type) {

	case *packet.ServerboundPlayConfirmTeleport:
		pl := s.getPlayerByConn(c)
		if pl == nil {
			panic("player not found")
		}
		return s.eventBus.Emit(&EventPlayerSpawn{
			Server: s,
			Player: pl,
		})
		return nil

	case *packet.ServerboundPlayUpdatePosition:
		pl := s.getPlayerByConn(c)
		if pl == nil {
			panic("player not found")
		}
		return s.eventBus.Emit(&EventPlayerMove{
			Server:      s,
			Player:      pl,
			X:           float64(p.X),
			Y:           float64(p.Y),
			Z:           float64(p.Z),
			HasPosition: true,
		})
		return nil

	case *packet.ServerboundPlayUpdatePositionAndRotation:
		pl := s.getPlayerByConn(c)
		if pl == nil {
			panic("player not found")
		}
		return s.eventBus.Emit(&EventPlayerMove{
			Server:      s,
			Player:      pl,
			X:           float64(p.X),
			Y:           float64(p.Y),
			Z:           float64(p.Z),
			HasPosition: true,
			Yaw:         float32(p.Yaw),
			Pitch:       float32(p.Pitch),
			HasRotation: true,
		})
		return nil

	case *packet.ServerboundPlayUpdateRotation:
		pl := s.getPlayerByConn(c)
		if pl == nil {
			panic("player not found")
		}
		return s.eventBus.Emit(&EventPlayerMove{
			Server:      s,
			Player:      pl,
			Yaw:         float32(p.Yaw),
			Pitch:       float32(p.Pitch),
			HasRotation: true,
		})
		return nil

	case *packet.ServerboundPlayChatMessage:
		pl := s.getPlayerByConn(c)
		if pl == nil {
			panic("player not found")
		}
		return s.eventBus.Emit(&EventPlayerChat{
			Server:  s,
			Player:  pl,
			Message: string(p.Message),
		})

	default:
		return nil
	}
}
