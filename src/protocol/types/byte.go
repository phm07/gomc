package types

import "io"

type Byte byte

func (b Byte) Marshal() []byte {
	return []byte{byte(b)}
}

func ReadByte(r io.Reader) (Byte, int, error) {
	b := make([]byte, 1)
	read, err := io.ReadFull(r, b)
	return Byte(b[0]), read, err
}
