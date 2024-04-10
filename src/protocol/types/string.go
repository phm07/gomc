package types

import (
	"io"
)

type String string

func (s String) Marshal() []byte {
	return append(VarInt(len(s)).Marshal(), []byte(s)...)
}

func ReadString(r io.Reader) (String, int, error) {
	length, n, err := ReadVarInt(r)
	if err != nil {
		return "", n, err
	}
	data := make([]byte, int(length))
	read, err := io.ReadFull(r, data)
	n += read
	if err != nil {
		return "", n, err
	}
	return String(data), n, nil
}
