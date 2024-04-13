package types

import (
	"gomc/src/util"
	"io"
	"math"
)

type Double float64

func (d Double) Marshal() []byte {
	return util.Uint64ToBytes(math.Float64bits(float64(d)))
}

func ReadDouble(r io.Reader) (Double, int, error) {
	var (
		b = make([]byte, 8)
		n int
	)
	read, err := io.ReadFull(r, b)
	n += read
	if err != nil {
		return 0, n, err
	}
	return Double(math.Float64frombits(util.Uint64FromBytes(b))), n, nil
}
