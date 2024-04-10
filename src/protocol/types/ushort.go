package types

import "io"

type UShort uint16

func (u UShort) Marshal() []byte {
	return []byte{byte(u >> 8), byte(u)}
}

func ReadUShort(r io.Reader) (UShort, int, error) {
	var (
		result UShort
		data   = make([]byte, 2)
		n      int
	)
	read, err := io.ReadFull(r, data)
	n += read
	if err != nil {
		return 0, n, err
	}
	result = UShort(data[0])<<8 | UShort(data[1])
	return result, n, nil
}
