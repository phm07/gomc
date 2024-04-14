package world

import (
	"log"
	"time"
)

type World struct {
	Height    int
	Chunks    map[int64]*Chunk
	Generator Generator
}

func NewWorld(height int, generator Generator) *World {
	return &World{
		Height:    height,
		Chunks:    make(map[int64]*Chunk),
		Generator: generator,
	}
}

func (w *World) GetOrGenerateChunk(x, z int) *Chunk {
	ch := w.Chunks[int64(x)<<32|(int64(z)&0xffffffff)]
	if ch == nil {
		start := time.Now()
		ch = w.Generator.Generate(w.Height, x, z)
		gen := time.Since(start)
		log.Printf("Generated chunk %d,%d in %s", x, z, gen)
		w.Chunks[int64(x)<<32|(int64(z)&0xffffffff)] = ch
	}
	return ch
}

func (w *World) GetHeightAt(x, z int) int {
	ch := w.GetOrGenerateChunk(x>>4, z>>4)
	return ch.GetHeightAt(x&0xf, z&0xf)
}
