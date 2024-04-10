package types

import "io"

type UUID []byte

func (u UUID) Marshal() []byte {
	return u
}

func ReadUUID(r io.Reader) (UUID, int, error) {
	u := make(UUID, 16)
	n, err := io.ReadFull(r, u)
	return u, n, err
}
