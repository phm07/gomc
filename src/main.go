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
	"gomc/src/encrypt"
	"gomc/src/protocol"
	"gomc/src/protocol/packet"
	"gomc/src/protocol/types"
	"gomc/src/registry"
	"gomc/src/session"
	"gomc/src/status"
	"io"
	"log"
	"net"
)

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
	}
	return nil
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
