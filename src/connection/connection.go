package connection

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"github.com/specspace/plasma/protocol/cfb8"
	"gomc/src/profile"
	"gomc/src/protocol"
	"gomc/src/protocol/packet"
	"gomc/src/protocol/types"
	"io"
	"net"
	"sync"
)

type Connection struct {
	io.Reader
	io.Writer
	Conn            net.Conn
	State           protocol.State
	ProtocolVersion int
	Username        string
	Closed          bool
	VerifyToken     []byte
	Secret          []byte
	Profile         *profile.Profile
	mu              *sync.Mutex
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{
		Conn:   conn,
		Reader: conn,
		Writer: conn,
		State:  protocol.StateHandshake,
		Closed: false,
		mu:     &sync.Mutex{},
	}
}

func (c *Connection) Close() error {
	c.Closed = true
	return c.Conn.Close()
}

func (c *Connection) SendPacket(s packet.SerializablePacket) error {
	if c.Closed {
		return nil
	}
	buf := bytes.NewBuffer(s.Serialize())
	id := types.VarInt(s.ID())
	length := types.VarInt(id.Len() + buf.Len())
	p := &protocol.Packet{
		Length:   length,
		PacketID: id,
		Data:     buf,
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	_, err := c.Write(p.Marshal())
	return err
}

func (c *Connection) ReadPacket() (p packet.SerializablePacket, read int, err error) {
	var (
		length, id types.VarInt
		n          int
	)

	length, n, err = types.ReadVarInt(c)
	read += n
	if err != nil {
		return
	}

	id, n, err = types.ReadVarInt(c)
	read += n
	if err != nil {
		return
	}

	buf := make([]byte, int(length)-id.Len())
	n, err = io.ReadFull(c, buf)
	if err != nil {
		return
	}

	p, err = packet.GetServerboundPacketInstance(c.State, int(id))
	if err != nil {
		return
	}

	err = p.Deserialize(buf)
	return
}

func (c *Connection) Encrypt() error {
	s, err := aes.NewCipher(c.Secret)
	if err != nil {
		return err
	}
	c.Reader = cipher.StreamReader{S: cfb8.NewDecrypter(s, c.Secret), R: c.Conn}
	c.Writer = cipher.StreamWriter{S: cfb8.NewEncrypter(s, c.Secret), W: c.Conn}
	return nil
}
