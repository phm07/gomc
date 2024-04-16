package server

import (
	"fmt"
	"gomc/src/protocol/packet"
	"gomc/src/protocol/types"
	"gomc/src/textcomponent"
	"slices"
	"time"
)

func (s *Server) registerListeners() {
	s.eventBus.RegisterListener(onPreLogin)
	s.eventBus.RegisterListener(onPlayerJoin)
	s.eventBus.RegisterListener(onPlayerQuit)
	s.eventBus.RegisterListener(onPlayerChat)
	s.eventBus.RegisterListener(onPlayerMove)
	s.eventBus.RegisterListener(onPlayerSpawn)
}

func onPreLogin(e *EventPreLogin) error {
	for _, p := range e.Server.players {
		if p.Name == e.Username || slices.Equal(p.UUID[:], e.UUID[:]) {
			e.Disallow(textcomponent.New("You are already connected"))
			break
		}
	}
	return nil
}

func onPlayerJoin(e *EventPlayerJoin) error {

	err := e.Player.Conn.SendPacket(&packet.ClientboundPlayLogin{
		EntityID:            0,
		IsHardcore:          false,
		DimensionNames:      []types.String{"world"},
		MaxPlayers:          0,
		ViewDistance:        types.VarInt(e.Server.cfg.ViewDistance),
		SimulationDistance:  8,
		ReducedDebugInfo:    false,
		EnableRespawnScreen: false,
		LimitedCrafting:     false,
		DimensionType:       "minecraft:overworld",
		DimensionName:       "world",
		HashedSeed:          0,
		GameMode:            0,
		PreviousGameMode:    -1,
		IsDebug:             false,
		IsFlat:              true,
		HasDeathLocation:    false,
		PortalCooldown:      0,
	})
	if err != nil {
		return err
	}

	h := e.Server.w.GetHeightAt(0, 0)
	err = e.Player.Conn.SendPacket(&packet.ClientboundPlaySynchronizePosition{
		X:          0,
		Y:          types.Double(float64(h) + 1.8),
		Z:          0,
		Yaw:        0,
		Pitch:      0,
		Flags:      0,
		TeleportID: 0,
	})
	if err != nil {
		return err
	}

	err = e.Player.Conn.SendPacket(&packet.ClientboundPlayPlayerCapabilities{
		Flags:       packet.PlayerCapabilityFlagAllowFlying,
		FlyingSpeed: 2.0,
	})
	if err != nil {
		return err
	}

	e.Server.BroadcastMessage(textcomponent.New(e.Player.Name + " joined the game").SetColor(textcomponent.ColorYellow))
	return nil
}

func onPlayerQuit(e *EventPlayerQuit) error {
	e.Server.BroadcastMessage(textcomponent.New(e.Player.Name + " left the game").SetColor(textcomponent.ColorYellow))
	return nil
}

func onPlayerChat(e *EventPlayerChat) error {
	e.Server.BroadcastMessage(textcomponent.New(fmt.Sprintf("<%s> %s", e.Player.Name, e.Message)))
	return nil
}

func onPlayerMove(e *EventPlayerMove) error {
	p := e.Player

	if e.HasRotation {
		p.Yaw, p.Pitch = e.Yaw, e.Pitch
	}

	if e.HasPosition {
		chunkX, chunkZ := int(e.X)>>4, int(e.Z)>>4
		prevChunkX, prevChunkZ := int(p.X)>>4, int(p.Z)>>4
		p.X, p.Y, p.Z = e.X, e.Y, e.Z

		if chunkX != prevChunkX || chunkZ != prevChunkZ {
			err := p.Conn.SendPacket(&packet.ClientboundPlaySetCenterChunk{
				ChunkX: types.VarInt(chunkX),
				ChunkZ: types.VarInt(chunkZ),
			})
			if err != nil {
				return err
			}
			vd := e.Server.cfg.ViewDistance
			for x := chunkX - vd; x <= chunkX+vd; x++ {
				for z := chunkZ - vd; z <= chunkZ+vd; z++ {
					if x < prevChunkX-vd || x > prevChunkX+vd || z < prevChunkZ-vd || z > prevChunkZ+vd {
						ch := e.Server.w.GetOrGenerateChunk(x, z)
						err = p.Conn.SendPacket(ch.Packet())
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}
	return nil
}

func onPlayerSpawn(e *EventPlayerSpawn) error {
	err := e.Player.Conn.SendPacket(&packet.ClientboundPlayPlayerInfoUpdate{
		Players: []*packet.ClientboundPlayPlayerInfoUpdatePlayer{
			{
				UUID: e.Player.Conn.Profile.Id[:],
				Actions: []packet.ClientboundPlayPlayerInfoUpdateAction{
					&packet.ClientboundPlayPlayerInfoUpdateActionAddPlayer{
						Name:       types.String(e.Player.Conn.Profile.Name),
						Properties: []*packet.ClientboundPlayPlayerInfoUpdateActionAddPlayerProperty{},
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}
	err = e.Player.Conn.SendPacket(&packet.ClientboundPlayGameEvent{
		Event: packet.ClientboundPlayGameEventWaitForChunks,
	})
	if err != nil {
		return err
	}
	err = e.Player.Conn.SendPacket(&packet.ClientboundPlaySetCenterChunk{
		ChunkX: 0,
		ChunkZ: 0,
	})
	if err != nil {
		return err
	}
	vd := e.Server.cfg.ViewDistance
	cch := e.Server.w.GetOrGenerateChunk(0, 0)
	err = e.Player.Conn.SendPacket(cch.Packet())
	if err != nil {
		return err
	}

	for x := -vd; x <= vd; x++ {
		for z := -vd; z <= vd; z++ {
			if x == 0 && z == 0 {
				continue
			}
			ch := e.Server.w.GetOrGenerateChunk(x, z)
			err = e.Player.Conn.SendPacket(ch.Packet())
			if err != nil {
				return err
			}
		}
	}

	go func() {
		t := time.NewTicker(10 * time.Second)
		for range t.C {
			if e.Player.Conn.Closed {
				return
			}
			_ = e.Player.Conn.SendPacket(&packet.ClientboundPlayKeepAlive{
				KeepAliveId: types.Long(time.Now().Unix()),
			})
		}
	}()
	return nil
}
