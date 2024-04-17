package world

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPalettedContainer_SetDataAt(t *testing.T) {
	data := make([]uint16, 4096)
	p := PalettedContainerFromBytes(data, 4, 8)
	for i := 0; i < len(data); i++ {
		assert.Equal(t, uint16(0), p.GetDataAt(i), i)
		p.SetDataAt(i, uint16(i))
		assert.Equal(t, uint16(i), p.GetDataAt(i), i)
	}
	for i := len(data) - 1; i >= 0; i-- {
		assert.Equal(t, uint16(i), p.GetDataAt(i), i)
		p.SetDataAt(i, uint16(len(data)-i-1))
		assert.Equal(t, uint16(len(data)-i-1), p.GetDataAt(i), i)
	}
	for i := 0; i < len(data); i++ {
		assert.Equal(t, uint16(len(data)-i-1), p.GetDataAt(i), i)
		p.SetDataAt(i, 0)
		assert.Equal(t, uint16(0), p.GetDataAt(i), i)
	}
	assert.Equal(t, p, &PalettedContainer{
		Size:      4096,
		Count:     map[uint16]int{0: 4096},
		BpeMin:    4,
		BpeThresh: 4,
		Palette:   &PaletteSingleValued{Value: 0},
	})
}

func TestGetBpe(t *testing.T) {
	testCases := map[int]int{
		2:  1,
		10: 4,
		16: 4,
		31: 5,
		32: 5,
	}
	for k, v := range testCases {
		assert.Equal(t, v, getBpe(k))
	}
}
