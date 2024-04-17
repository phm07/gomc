package world

import (
	"math/rand"
)

type Generator interface {
	Generate(w *World, x, z int) *Chunk
}

type RandomGenerator struct {
	Blocks []uint16
	Height int
}

func (g *RandomGenerator) Generate(w *World, x, z int) *Chunk {
	src := rand.NewSource(int64(x)<<32 | (int64(z) & 0xffffffff))
	rng := rand.New(src)
	data := make([]uint16, w.Height<<8)
	for i := 0; i < g.Height<<8; i++ {
		data[i] = g.Blocks[rng.Intn(len(g.Blocks))]
	}
	c := NewChunk(w, x, z, data)
	c.CalculateSkyLight()
	return c
}
