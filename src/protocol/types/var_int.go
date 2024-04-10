package types

import (
	"bytes"
	"io"
)

type VarInt int32

func (v VarInt) Marshal() []byte {
	var buf bytes.Buffer
	u := uint32(v)
	for {
		temp := (byte)(u & 0x7f)
		u >>= 7
		if u != 0 {
			temp |= 0x80
		}
		buf.WriteByte(temp)
		if u == 0 {
			break
		}
	}
	return buf.Bytes()
}

func ReadVarInt(r io.Reader) (VarInt, int, error) {
	var (
		result VarInt
		shift  uint
		n      int
		b      = make([]byte, 1)
	)
	for {
		read, err := io.ReadFull(r, b)
		n += read
		if err != nil {
			return 0, n, err
		}
		value := VarInt(b[0] & 0x7f)
		result |= value << shift
		shift += 7
		if b[0]&0x80 == 0 {
			break
		}
	}
	return result, n, nil
}

func (v VarInt) Len() int {
	var n int
	u := uint32(v)
	for {
		n++
		u >>= 7
		if u == 0 {
			break
		}
	}
	return n
}
