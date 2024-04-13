package world

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPack(t *testing.T) {
	testData := []struct {
		data     []uint16
		expected []uint64
		bpe      int
	}{
		{
			data:     []uint16{1, 2, 2, 3, 4, 4, 5, 6, 6, 4, 8, 0, 7, 4, 3, 13, 15, 16, 9, 14, 10, 12, 0, 2},
			expected: []uint64{0x0020863148418841, 0x01018A7260F68C87},
			bpe:      5,
		},
		{
			data:     []uint16{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			expected: []uint64{0xfedcba987654321},
			bpe:      4,
		},
	}

	for _, test := range testData {
		packed := pack(test.data, test.bpe)
		assert.Equal(t, test.expected, packed)
	}
}
