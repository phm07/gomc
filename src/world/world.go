package world

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
		ch = w.Generator.Generate(w.Height, x, z)
		w.Chunks[int64(x)<<32|(int64(z)&0xffffffff)] = ch
	}
	return ch
}
