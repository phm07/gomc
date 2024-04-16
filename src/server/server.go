package server

import (
	"errors"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/spf13/viper"
	"gomc/src/connection"
	"gomc/src/event"
	"gomc/src/player"
	"gomc/src/protocol"
	"gomc/src/protocol/packet"
	"gomc/src/world"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"
)

type Server struct {
	cfg       *Config
	w         *world.World
	listener  net.Listener
	stopCh    chan struct{}
	players   []*player.Player
	playersMu *sync.Mutex
	eventBus  *event.Bus
}

func NewServer(cfg *Config) *Server {
	srv := &Server{
		cfg:       cfg,
		w:         world.NewWorld(384, 4, &world.NaturalGenerator{}),
		stopCh:    make(chan struct{}),
		playersMu: &sync.Mutex{},
		eventBus:  event.NewBus(),
	}
	srv.registerListeners()
	return srv
}

func (s *Server) Start() {
	vd := s.cfg.ViewDistance

	log.Println("Generating chunks...")
	start := time.Now()
	for x := -vd; x <= vd; x++ {
		for z := -vd; z <= vd; z++ {
			_ = s.w.GetOrGenerateChunk(x, z)
		}
	}

	t := time.Since(start)
	nChunks := float64(vd*2+1) * float64(vd*2+1)
	log.Printf("Generated %.0f chunks in %s (%.2f cps, avg %s)\n",
		nChunks, t, nChunks/t.Seconds(), time.Duration(int64(t)/int64(nChunks)))

	log.Printf("World size: %s", humanize.Bytes(uint64(s.w.Size())))
	go func() {
		t := time.NewTicker(10 * time.Second)
		for range t.C {
			log.Printf("World size: %s", humanize.Bytes(uint64(s.w.Size())))
		}
	}()

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)
		<-ch
		close(s.stopCh)
	}()

	go s.listen()
	<-s.stopCh

	log.Println("Shutting down...")

	if s.listener != nil {
		_ = s.listener.Close()
	}
}

func (s *Server) listen() {
	addr := viper.GetString("bind_addr")
	port := viper.GetInt("port")

	var err error
	s.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		panic(err)
	}

	log.Printf("Listening on %s:%d\n", addr, port)

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			break
		}

		go s.handleConnection(connection.NewConnection(conn))
	}
}

func (s *Server) handleConnection(c *connection.Connection) {
	/*defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v\n", r)
		}
	}()*/

	log.Printf("Client %s connected\n", c.Conn.RemoteAddr())

	defer func() {
		log.Printf("Client %s disconnected\n", c.Conn.RemoteAddr())
		pl := s.getPlayerByConn(c)
		if pl != nil {
			_ = s.eventBus.Emit(&EventPlayerQuit{
				Server: s,
				Player: pl,
			})

			s.playersMu.Lock()
			defer s.playersMu.Unlock()
			for i, p := range s.players {
				if p == pl {
					s.players = append(s.players[:i], s.players[i+1:]...)
					break
				}
			}
		}
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

		if err := s.handlePacket(c, p); err != nil {
			log.Printf("error while handling packet: %s\n", err)
			continue
		}
	}
}

func (s *Server) handlePacket(c *connection.Connection, p packet.SerializablePacket) error {

	switch c.State {
	case protocol.StateHandshake:
		return s.handleHandshake(c, p)
	case protocol.StateStatus:
		return s.handleStatus(c, p)
	case protocol.StateConfiguration:
		return s.handleConfiguration(c, p)
	case protocol.StateLogin:
		return s.handleLogin(c, p)
	case protocol.StatePlay:
		return s.handlePlay(c, p)
	default:
		return nil
	}
}

func (s *Server) getPlayerByConn(c *connection.Connection) *player.Player {
	s.playersMu.Lock()
	defer s.playersMu.Unlock()
	for _, p := range s.players {
		if p.Conn == c {
			return p
		}
	}
	return nil
}
