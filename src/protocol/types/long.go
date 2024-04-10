package types

import "io"

type Long int64

func (l Long) Marshal() []byte {
	b := make([]byte, 8)
	for i := 0; i < 8; i++ {
		b[i] = byte(l >> (i * 8))
	}
	return b
}

func ReadLong(r io.Reader) (Long, int, error) {
	var (
		l Long
		b = make([]byte, 8)
		n int
	)
	read, err := io.ReadFull(r, b)
	n += read
	if err != nil {
		return 0, n, err
	}
	for i := 0; i < 8; i++ {
		l |= Long(b[i]) << (i * 8)
	}
	return l, n, nil
}
