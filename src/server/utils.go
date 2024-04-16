package server

import (
	"gomc/src/textcomponent"
	"log"
)

func (s *Server) BroadcastMessage(msg *textcomponent.Component) {
	for _, p := range s.players {
		_ = p.SendMessage(msg)
	}
	log.Println(msg.Plain())
}
