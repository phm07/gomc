package types

import (
	"gomc/src/util"
	"io"
)

type Int int32

func (v Int) Marshal() []byte {
	return util.Int32ToBytes(int32(v))
}

func ReadInt(r io.Reader) (Int, int, error) {
	var (
		b = make([]byte, 4)
		n int
	)
	read, err := io.ReadFull(r, b)
	n += read
	if err != nil {
		return 0, n, err
	}
	return Int(util.Int32FromBytes(b)), n, nil
}
