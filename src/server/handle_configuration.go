package server

import (
	"gomc/src/connection"
	"gomc/src/protocol"
	"gomc/src/protocol/packet"
	"gomc/src/registry"
)

func (s *Server) handleConfiguration(c *connection.Connection, p packet.SerializablePacket) error {
	switch p.(type) {

	case *packet.ServerboundConfigurationClientInformation:
		err := c.SendPacket(&packet.ClientboundConfigurationRegistryData{RegistryDataNBT: registry.RegistryNBTBytes})
		if err != nil {
			return err
		}
		return c.SendPacket(&packet.ClientboundConfigurationFinish{})

	case *packet.ServerboundConfigurationFinishAck:
		c.State = protocol.StatePlay
		pl := s.getPlayerByConn(c)
		if pl == nil {
			panic("player not found")
		}
		return s.eventBus.Emit(&EventPlayerJoin{
			Server: s,
			Player: pl,
		})

	default:
		return nil
	}
}
