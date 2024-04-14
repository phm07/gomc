package types

import (
	"bytes"
	"gomc/src/util"
	"io"
)

type BitSet []uint64

func (b BitSet) Marshal() []byte {
	res := make([]byte, len(b)<<3)
	for i, v := range b {
		l := util.Uint64ToBytes(v)
		for j := 0; j < 8; j++ {
			res[(i<<3)+j] = l[j]
		}
	}
	var buf bytes.Buffer
	buf.Write(VarInt(len(b)).Marshal())
	buf.Write(res)
	return buf.Bytes()
}

func ReadBitSet(r io.Reader) (BitSet, int, error) {
	length, n, err := ReadVarInt(r)
	if err != nil {
		return nil, n, err
	}
	data := make([]byte, int(length)<<3)
	read, err := io.ReadFull(r, data)
	n += read
	if err != nil {
		return nil, n, err
	}
	res := make(BitSet, length)
	for i := 0; i < len(res); i++ {
		res[i] = util.Uint64FromBytes(data[(i << 3):((i + 1) << 3)])
	}
	return res, n, nil
}

func NewBitSet(bits int) BitSet {
	return make(BitSet, (bits+63)>>6)
}

func (b BitSet) SetBit(index int, value bool) {
	if value {
		b[index>>6] |= 1 << (index & 0x3f)
	} else {
		b[index>>6] &^= 1 << (index & 0x3f)
	}
}

func (b BitSet) GetBit(index int) bool {
	return b[index>>6]&(1<<(index&0x3f)) != 0
}
