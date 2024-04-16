package types

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPosition_Marshal(t *testing.T) {
	p := Position{X: 18357644, Y: 831, Z: -20882616}
	assert.Equal(t, []byte{0x46, 0x7, 0x63, 0x2c, 0x15, 0xb4, 0x83, 0x3f}, p.Marshal())
}

func TestReadPosition(t *testing.T) {
	b := []byte{0x46, 0x7, 0x63, 0x2c, 0x15, 0xb4, 0x83, 0x3f}
	p, n, err := ReadPosition(bytes.NewReader(b))
	assert.NoError(t, err)
	assert.Equal(t, 8, n)
	assert.Equal(t, Position{X: 18357644, Y: 831, Z: -20882616}, p)
}
