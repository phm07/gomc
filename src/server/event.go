package server

import (
	"github.com/google/uuid"
	"gomc/src/connection"
	"gomc/src/player"
	"gomc/src/textcomponent"
	"gomc/src/world"
)

type EventPreLogin struct {
	Server   *Server
	Conn     *connection.Connection
	UUID     uuid.UUID
	Username string
	Disallow func(reason *textcomponent.Component)
}

type EventPlayerLogin struct {
	Server *Server
	Player *player.Player
}

type EventPlayerJoin struct {
	Server *Server
	Player *player.Player
}

type EventPlayerQuit struct {
	Server *Server
	Player *player.Player
}

type EventPlayerChat struct {
	Server  *Server
	Player  *player.Player
	Message string
}

type EventPlayerMove struct {
	Server      *Server
	Player      *player.Player
	X, Y, Z     float64
	HasPosition bool
	Yaw, Pitch  float32
	HasRotation bool
}

type EventPlayerSpawn struct {
	Server *Server
	Player *player.Player
}

type EventPlayerBlockBreak struct {
	Server *Server
	Player *player.Player
	Block  *world.Block
}
