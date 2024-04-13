package types

import (
	"gomc/src/util"
	"io"
	"math"
)

type Float float32

func (d Float) Marshal() []byte {
	return util.Uint32ToBytes(math.Float32bits(float32(d)))
}

func ReadFloat(r io.Reader) (Float, int, error) {
	var (
		b = make([]byte, 4)
		n int
	)
	read, err := io.ReadFull(r, b)
	n += read
	if err != nil {
		return 0, n, err
	}
	return Float(math.Float32frombits(util.Uint32FromBytes(b))), n, nil
}
