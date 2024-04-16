package server

import (
	"fmt"
	"gomc/src/connection"
	"gomc/src/protocol"
	"gomc/src/protocol/packet"
	"log"
)

func (s *Server) handleHandshake(c *connection.Connection, p packet.SerializablePacket) error {
	switch p := p.(type) {
	case *packet.ServerboundHandshake:
		log.Printf("Protocol version: %d\n", p.ProtocolVersion)
		log.Printf("Address: %s\n", p.ServerAddress)
		log.Printf("Port: %d\n", p.ServerPort)
		log.Printf("Next state: %d\n", p.NextState)

		c.ProtocolVersion = int(p.ProtocolVersion)
		switch p.NextState {
		case 1:
			c.State = protocol.StateStatus
		case 2:
			c.State = protocol.StateLogin
		default:
			return fmt.Errorf("unknown next state: %d", p.NextState)
		}
		return nil

	default:
		return nil
	}
}
