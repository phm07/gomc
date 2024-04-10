package packet

import (
	"fmt"
	"gomc/src/protocol"
)

type SerializablePacket interface {
	Serialize() []byte
	Deserialize([]byte) error
	State() protocol.State
	ID() int
}

var serverbound = make(map[protocol.State]map[int]func() SerializablePacket)

func init() {
	for i := 0; i < 4; i++ {
		serverbound[protocol.State(i)] = make(map[int]func() SerializablePacket)
	}
}

func GetServerboundPacketInstance(state protocol.State, id int) (SerializablePacket, error) {
	m, ok := serverbound[state]
	if !ok {
		return nil, fmt.Errorf("unknown state %d", state)
	}
	c, ok := m[id]
	if !ok {
		return nil, fmt.Errorf("unknown id %d", id)
	}
	return c(), nil
}
