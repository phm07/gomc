package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLog2(t *testing.T) {
	tests := map[int]int{
		0:   0,
		1:   0,
		2:   1,
		3:   1,
		4:   2,
		5:   2,
		6:   2,
		7:   2,
		8:   3,
		15:  3,
		16:  4,
		32:  5,
		64:  6,
		128: 7,
	}
	for n, expected := range tests {
		assert.Equal(t, expected, Log2(n))
	}
}
