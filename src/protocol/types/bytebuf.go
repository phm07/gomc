package types

import (
	"bytes"
	"io"
)

type ByteBuf []byte

func (b ByteBuf) Marshal() []byte {
	var buf bytes.Buffer
	buf.Write(VarInt(len(b)).Marshal())
	buf.Write(b)
	return buf.Bytes()
}

func ReadByteBuf(r io.Reader) (ByteBuf, int, error) {
	read := 0
	length, n, err := ReadVarInt(r)
	read += n
	if err != nil {
		return nil, read, err
	}
	buf := make([]byte, length)
	n, err = io.ReadFull(r, buf)
	read += n
	if err != nil {
		return nil, read, err
	}
	return buf, read, nil
}
