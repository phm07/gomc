package types

import (
	"gomc/src/util"
	"io"
)

type UShort uint16

func (u UShort) Marshal() []byte {
	return util.Uint16ToBytes(uint16(u))
}

func ReadUShort(r io.Reader) (UShort, int, error) {
	var (
		b = make([]byte, 2)
		n int
	)
	read, err := io.ReadFull(r, b)
	n += read
	if err != nil {
		return 0, n, err
	}
	return UShort(util.Uint16FromBytes(b)), n, nil
}
