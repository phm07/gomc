package types_test

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"gomc/src/protocol/types"
	"strconv"
	"testing"
)

var varIntTests = []struct {
	value int
	bytes []byte
}{
	{0, []byte{0x00}},
	{1, []byte{0x01}},
	{2, []byte{0x02}},
	{127, []byte{0x7f}},
	{128, []byte{0x80, 0x01}},
	{255, []byte{0xff, 0x01}},
	{25565, []byte{0xdd, 0xc7, 0x01}},
	{2097151, []byte{0xff, 0xff, 0x7f}},
	{2147483647, []byte{0xff, 0xff, 0xff, 0xff, 0x07}},
	{-1, []byte{0xff, 0xff, 0xff, 0xff, 0x0f}},
	{-2147483648, []byte{0x80, 0x80, 0x80, 0x80, 0x08}},
}

func TestReadVarInt(t *testing.T) {
	for _, test := range varIntTests {
		t.Run(strconv.Itoa(test.value), func(t *testing.T) {
			value, n, err := types.ReadVarInt(bytes.NewReader(test.bytes))
			assert.Equal(t, types.VarInt(test.value), value)
			assert.Equal(t, len(test.bytes), n)
			assert.NoError(t, err)
		})
	}
}

func TestWriteVarInt(t *testing.T) {
	for _, test := range varIntTests {
		t.Run(strconv.Itoa(test.value), func(t *testing.T) {
			b := types.VarInt(test.value).Marshal()
			assert.Equal(t, test.bytes, b)
		})
	}
}

func TestVarInt_Len(t *testing.T) {
	for _, test := range varIntTests {
		t.Run(strconv.Itoa(test.value), func(t *testing.T) {
			assert.Equal(t, len(test.bytes), types.VarInt(test.value).Len())
		})
	}
}
