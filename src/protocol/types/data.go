package types

import (
	"io"
)

type Data []byte

func (d Data) Marshal() []byte {
	return d
}

func ReadData(r io.Reader) (Data, int, error) {
	d, err := io.ReadAll(r)
	return d, len(d), err
}
