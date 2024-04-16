package world

import (
	"math/rand"
)

type Generator interface {
	Generate(seed int64, height, x, z int) *Chunk
}

type RandomGenerator struct {
	Blocks []uint16
	Height int
}

func (g *RandomGenerator) Generate(height, x, z int) *Chunk {
	c := NewChunk(height, x, z)
	src := rand.NewSource(int64(x)<<32 | (int64(z) & 0xffffffff))
	rng := rand.New(src)
	for i := 0; i < g.Height<<8; i++ {
		c.Data[i] = g.Blocks[rng.Intn(len(g.Blocks))]
	}
	c.CalculateSkyLight()
	return c
}
