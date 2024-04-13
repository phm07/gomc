package types

import (
	"gomc/src/util"
	"io"
)

type Long int64

func (l Long) Marshal() []byte {
	return util.Int64ToBytes(int64(l))
}

func ReadLong(r io.Reader) (Long, int, error) {
	var (
		b = make([]byte, 8)
		n int
	)
	read, err := io.ReadFull(r, b)
	n += read
	if err != nil {
		return 0, n, err
	}
	return Long(util.Int64FromBytes(b)), n, nil
}
