package world

type World struct {
	Height    int
	Chunks    map[int64]*Chunk
	Generator Generator
	Seed      int64
}

func NewWorld(height int, seed int64, generator Generator) *World {
	return &World{
		Height:    height,
		Chunks:    make(map[int64]*Chunk),
		Generator: generator,
		Seed:      seed,
	}
}

func (w *World) Size() uint64 {
	return uint64(len(w.Chunks) * ((16 * 16 * w.Height * 5) >> 1))
}

func (w *World) GetOrGenerateChunk(x, z int) *Chunk {
	ch := w.Chunks[int64(x)<<32|(int64(z)&0xffffffff)]
	if ch == nil {
		ch = w.Generator.Generate(w.Seed, w.Height, x, z)
		w.Chunks[int64(x)<<32|(int64(z)&0xffffffff)] = ch
	}
	return ch
}

func (w *World) GetHeightAt(x, z int) int {
	ch := w.GetOrGenerateChunk(x>>4, z>>4)
	return ch.GetHeightAt(x&0xf, z&0xf)
}
