package types

import "io"

type Boolean bool

func (v Boolean) Marshal() []byte {
	var b byte = 0
	if v {
		b = 1
	}
	return []byte{b}
}

func ReadBoolean(r io.Reader) (Boolean, int, error) {
	b := make([]byte, 1)
	read, err := io.ReadFull(r, b)
	return b[0] != 0, read, err
}
