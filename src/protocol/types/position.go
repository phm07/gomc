package types

import (
	"gomc/src/util"
	"io"
)

type Position struct {
	X, Y, Z int
}

func (p Position) Marshal() []byte {
	var l uint64
	l |= uint64(p.Y) & 0xfff
	l |= (uint64(p.Z) & 0x3ffffff) << 12
	l |= (uint64(p.X) & 0x3ffffff) << 38
	return util.Uint64ToBytes(l)
}

func ReadPosition(r io.Reader) (Position, int, error) {
	var (
		b = make([]byte, 8)
		n int
	)
	read, err := io.ReadFull(r, b)
	n += read
	if err != nil {
		return Position{}, n, err
	}
	l := util.Uint64FromBytes(b)
	p := Position{}
	p.Y = int(int64(l<<52) >> 52)
	p.Z = int(int64(l<<26) >> 38)
	p.X = int(int64(l) >> 38)
	return p, n, err
}
