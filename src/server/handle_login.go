package server

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/google/uuid"
	"gomc/src/connection"
	"gomc/src/encrypt"
	"gomc/src/player"
	"gomc/src/profile"
	"gomc/src/protocol"
	"gomc/src/protocol/packet"
	"gomc/src/protocol/types"
	"gomc/src/session"
	"gomc/src/textcomponent"
	"io"
	"log"
)

func (s *Server) handleLogin(c *connection.Connection, p packet.SerializablePacket) error {
	switch p := p.(type) {

	case *packet.ServerboundLoginStart:
		c.Username = string(p.Username)

		var (
			disallowed bool
			reason     *textcomponent.Component
		)
		err := s.eventBus.Emit(&EventPreLogin{
			Server:   s,
			Conn:     c,
			UUID:     uuid.UUID(p.UUID),
			Username: string(p.Username),
			Disallow: func(r *textcomponent.Component) {
				disallowed = true
				reason = r
			},
		})
		if err != nil {
			return err
		}

		if disallowed {
			err = c.SendPacket(&packet.ClientboundLoginDisconnect{
				Reason: types.String(reason.MarshalJSON()),
			})
			if err != nil {
				return err
			}
			return c.Close()
		}

		if s.cfg.OnlineMode {
			c.VerifyToken = make([]byte, 4)
			_, err = io.ReadFull(rand.Reader, c.VerifyToken)
			if err != nil {
				return err
			}
			return c.SendPacket(&packet.ClientboundLoginEncryptionRequest{
				PublicKey:   encrypt.PublicKeyBytes,
				VerifyToken: c.VerifyToken,
			})

		} else {
			c.Profile = profile.OfflineProfile(string(p.Username))
			return c.SendPacket(&packet.ClientboundLoginSuccess{
				UUID:     c.Profile.Id[:],
				Username: types.String(c.Profile.Name),
			})
		}

	case *packet.ServerboundLoginEncryptionResponse:
		if !s.cfg.OnlineMode {
			return fmt.Errorf("server is in offline mode")
		}

		verify, err := encrypt.Decrypt(p.VerifyToken)
		if err != nil {
			return err
		}
		if !bytes.Equal(verify, c.VerifyToken) {
			return fmt.Errorf("invalid verify token")
		}

		c.Secret, err = encrypt.Decrypt(p.SharedSecret)
		if err != nil {
			return err
		}
		err = c.Encrypt()
		if err != nil {
			return err
		}

		log.Printf("Connection to %s is now encrypted", c.Conn.RemoteAddr())

		c.Profile, err = session.ValidateSession(c)
		if err != nil {
			log.Printf("Failed to validate session: %v\n", err)
			return c.Close()
		}

		return c.SendPacket(&packet.ClientboundLoginSuccess{
			Username: types.String(c.Profile.Name),
			UUID:     c.Profile.Id[:],
		})

	case *packet.ServerboundLoginAck:
		c.State = protocol.StateConfiguration
		pl := &player.Player{
			Conn: c,
			UUID: c.Profile.Id,
			Name: c.Profile.Name,
			EID:  s.nextEid(),
		}
		s.playersMu.Lock()
		s.players = append(s.players, pl)
		s.playersMu.Unlock()
		return s.eventBus.Emit(&EventPlayerLogin{
			Server: s,
			Player: pl,
		})

	default:
		return nil
	}
}
