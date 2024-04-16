package server

import (
	"errors"
	"fmt"
	"gomc/src/connection"
	"gomc/src/protocol/packet"
	"gomc/src/protocol/types"
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

	case *packet.ServerboundPlayPlayerAction:
		switch p.Status {
		case packet.PlayerActionStatusFinishDigging:
			if p.Location.Y < 0 || p.Location.Y > s.w.Height {
				return fmt.Errorf("invalid y: %d", p.Location.Y)
			}
			b := s.w.BlockAt(p.Location.X, p.Location.Y, p.Location.Z)
			e := &EventPlayerBlockBreak{
				Server: s,
				Player: s.getPlayerByConn(c),
				Block:  b,
			}
			err := s.eventBus.Emit(e)
			err = errors.Join(err, c.SendPacket(&packet.ClientboundPlayBlockUpdate{
				Location: p.Location,
				BlockID:  types.VarInt(b.GetState()),
			}))
			return errors.Join(err, c.SendPacket(&packet.ClientboundPlayAckBlockChange{
				SequenceID: p.SequenceID,
			}))
		}
		return nil

	default:
		return nil
	}
}
