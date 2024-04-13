package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"github.com/ztrue/tracerr"
	"gomc/src/connection"
	"gomc/src/data"
	"gomc/src/encrypt"
	"gomc/src/protocol"
	"gomc/src/protocol/packet"
	"gomc/src/protocol/types"
	"gomc/src/registry"
	"gomc/src/session"
	"gomc/src/status"
	"gomc/src/world"
	"io"
	"log"
	"net"
	"time"
)

var w = world.NewWorld(384, &world.RandomGenerator{
	Blocks: []data.Block{data.DiamondOre, data.CoalOre, data.RedstoneOre, data.Stone, data.Dirt, data.DiamondBlock},
	Height: 128,
})

const viewDistance = 3

func main() {
	if err := loadConfig(); err != nil {
		panic(err)
	}

	if err := status.Init(); err != nil {
		panic(err)
	}

	if err := encrypt.GenerateKeypair(); err != nil {
		panic(err)
	}

	addr := viper.GetString("bind_addr")
	port := viper.GetInt("port")

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = listener.Close()
	}()

	log.Printf("Listening on %s:%d\n", addr, port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			tracerr.PrintSourceColor(err)
			continue
		}

		go handleRequest(connection.NewConnection(conn))
	}
}

func handleRequest(c *connection.Connection) {
	log.Printf("Client %s connected\n", c.Conn.RemoteAddr())

	defer func() {
		log.Printf("Client %s disconnected\n", c.Conn.RemoteAddr())
	}()

	for !c.Closed {
		p, _, err := c.ReadPacket()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			log.Printf("error while reading packet: %s\n", err)
			continue
		}

		if err := handlePacket(c, p); err != nil {
			tracerr.PrintSourceColor(err)
			continue
		}
	}
}

func handlePacket(c *connection.Connection, p packet.SerializablePacket) error {
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

	case *packet.ServerboundStatusRequest:
		stat := status.GetStatus()
		stat.Version.Protocol = c.ProtocolVersion
		statBytes, err := json.Marshal(stat)
		if err != nil {
			return err
		}

		return c.SendPacket(&packet.ClientboundStatusResponse{
			Json: types.String(statBytes),
		})

	case *packet.ServerboundStatusPing:
		return c.SendPacket(&packet.ClientboundStatusPong{Payload: p.Payload})

	case *packet.ServerboundLoginStart:
		c.Username = string(p.Username)
		c.VerifyToken = make([]byte, 4)
		_, err := io.ReadFull(rand.Reader, c.VerifyToken)
		if err != nil {
			return tracerr.Wrap(err)
		}
		return c.SendPacket(&packet.ClientboundLoginEncryptionRequest{
			PublicKey:   encrypt.PublicKeyBytes,
			VerifyToken: c.VerifyToken,
		})

	case *packet.ServerboundLoginEncryptionResponse:
		verify, err := encrypt.Decrypt(p.VerifyToken)
		if err != nil {
			return tracerr.Wrap(err)
		}
		if !bytes.Equal(verify, c.VerifyToken) {
			return fmt.Errorf("invalid verify token")
		}

		c.Secret, err = encrypt.Decrypt(p.SharedSecret)
		if err != nil {
			return tracerr.Wrap(err)
		}
		err = c.Encrypt()
		if err != nil {
			return tracerr.Wrap(err)
		}

		log.Printf("Connection to %s is now encrypted", c.Conn.RemoteAddr())

		err = session.ValidateSession(c)
		if err != nil {
			log.Printf("Failed to validate session: %v\n", err)
			return c.Close()
		}

		err = c.SendPacket(&packet.ClientboundLoginSuccess{
			Username: types.String(c.Profile.Name),
			UUID:     c.Profile.Id[:],
		})

	case *packet.ServerboundLoginAck:
		c.State = protocol.StateConfiguration

	case *packet.ServerboundConfigurationClientInformation:
		err := c.SendPacket(&packet.ClientboundConfigurationRegistryData{RegistryDataNBT: registry.RegistryNBTBytes})
		if err != nil {
			return err
		}
		return c.SendPacket(&packet.ClientboundConfigurationFinish{})

	case *packet.ServerboundConfigurationPluginMessage:
		// idc

	case *packet.ServerboundConfigurationFinishAck:
		c.State = protocol.StatePlay

		err := c.SendPacket(&packet.ClientboundPlayLogin{
			EntityID:            0,
			IsHardcore:          false,
			DimensionNames:      []types.String{"world"},
			MaxPlayers:          0,
			ViewDistance:        8,
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
		return c.SendPacket(&packet.ClientboundPlaySynchronizePosition{
			X:          0,
			Y:          66,
			Z:          0,
			Yaw:        0,
			Pitch:      0,
			Flags:      0,
			TeleportID: 0,
		})

	case *packet.ServerboundConfirmTeleport:
		err := c.SendPacket(&packet.ClientboundPlayPlayerInfoUpdate{
			Players: []*packet.ClientboundPlayPlayerInfoUpdatePlayer{
				{
					UUID: c.Profile.Id[:],
					Actions: []packet.ClientboundPlayPlayerInfoUpdateAction{
						&packet.ClientboundPlayPlayerInfoUpdateActionAddPlayer{
							Name:       types.String(c.Profile.Name),
							Properties: []*packet.ClientboundPlayPlayerInfoUpdateActionAddPlayerProperty{},
						},
					},
				},
			},
		})
		if err != nil {
			return err
		}
		err = c.SendPacket(&packet.ClientboundPlayGameEvent{
			Event: packet.ClientboundPlayGameEventWaitForChunks,
		})
		if err != nil {
			return err
		}
		err = c.SendPacket(&packet.ClientboundPlaySetCenterChunk{
			ChunkX: 0,
			ChunkZ: 0,
		})
		if err != nil {
			return err
		}

		for x := -viewDistance; x <= viewDistance; x++ {
			for z := -viewDistance; z <= viewDistance; z++ {
				ch := w.GetOrGenerateChunk(x, z)
				err = sendChunk(c, ch)
				if err != nil {
					return err
				}
			}
		}

		go func() {
			t := time.NewTicker(10 * time.Second)
			for range t.C {
				if c.Closed {
					return
				}
				_ = c.SendPacket(&packet.ClientboundPlayKeepAlive{
					KeepAliveId: types.Long(time.Now().Unix()),
				})
			}
		}()

	case *packet.ServerboundPlayKeepAlive:
		return c.SendPacket(&packet.ClientboundPlayKeepAlive{
			KeepAliveId: p.KeepAliveId,
		})

	case *packet.ServerboundPlayUpdatePosition:
		chunkX, chunkZ := int(p.X)>>4, int(p.Z)>>4
		prevChunkX, prevChunkZ := int(c.X)>>4, int(c.Z)>>4
		c.X, c.Y, c.Z = float64(p.X), float64(p.Y), float64(p.Z)
		if chunkX != prevChunkX || chunkZ != prevChunkZ {
			err := c.SendPacket(&packet.ClientboundPlaySetCenterChunk{
				ChunkX: types.VarInt(chunkX),
				ChunkZ: types.VarInt(chunkZ),
			})
			if err != nil {
				return err
			}
			for x := chunkX - viewDistance; x <= chunkX+viewDistance; x++ {
				for z := chunkZ - viewDistance; z <= chunkZ+viewDistance; z++ {
					if x < prevChunkX-viewDistance || x > prevChunkX+viewDistance || z < prevChunkZ-viewDistance || z > prevChunkZ+viewDistance {
						ch := w.GetOrGenerateChunk(x, z)
						err = sendChunk(c, ch)
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

func sendChunk(c *connection.Connection, ch *world.Chunk) error {
	mask := types.NewBitSet(26)
	mask2 := types.NewBitSet(26)
	for i := 0; i < 26; i++ {
		mask2.SetBit(i, true)
	}
	return c.SendPacket(&packet.ClientboundPlayChunkData{
		ChunkX:               types.Int(ch.X),
		ChunkZ:               types.Int(ch.Z),
		Heightmaps:           ch.HeightMap().Marshal(),
		Data:                 ch.Marshal(),
		NumBlockEntities:     0,
		SkyLightMask:         mask,
		BlockLightMask:       mask,
		EmptySkyLightMask:    mask2,
		EmptyBlockLightMask:  mask2,
		SkyLightArrayCount:   0,
		BlockLightArrayCount: 0,
	})
}

func loadConfig() error {
	viper.SetDefault("bind_addr", "")
	viper.SetDefault("port", 25565)
	viper.SetDefault("motd", "Hello world!")
	viper.SetDefault("max_players", 100)
	viper.SetDefault("online_mode", true)

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	var configFileNotFoundError viper.ConfigFileNotFoundError
	if errors.As(err, &configFileNotFoundError) {
		return viper.SafeWriteConfigAs("config.toml")
	}
	return err
}
