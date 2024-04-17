package world

import "sync"

type World struct {
	Height     int
	MinY, MaxY int
	Chunks     map[int64]*Chunk
	Generator  Generator
	Seed       int64
	chunksMu   *sync.RWMutex
}

func NewWorld(height int, minY int, seed int64, generator Generator) *World {
	return &World{
		Height:    height,
		MinY:      minY,
		MaxY:      height + minY - 1,
		Chunks:    make(map[int64]*Chunk),
		Generator: generator,
		Seed:      seed,
		chunksMu:  &sync.RWMutex{},
	}
}

func (w *World) Size() uint64 {
	w.chunksMu.RLock()
	defer w.chunksMu.RUnlock()
	var size uint64
	for _, c := range w.Chunks {
		size += uint64(len(c.Marshal()))
	}
	return size
}

func (w *World) GetOrGenerateChunk(x, z int) *Chunk {
	w.chunksMu.RLock()
	ch := w.Chunks[int64(x)<<32|(int64(z)&0xffffffff)]
	w.chunksMu.RUnlock()
	if ch == nil {
		ch = w.Generator.Generate(w, x, z)
		w.chunksMu.Lock()
		w.Chunks[int64(x)<<32|(int64(z)&0xffffffff)] = ch
		w.chunksMu.Unlock()
	}
	return ch
}

func (w *World) GetHeightAt(x, z int) int {
	ch := w.GetOrGenerateChunk(x>>4, z>>4)
	return ch.GetHeightAt(x&0xf, z&0xf)
}

type Block struct {
	World   *World
	Chunk   *Chunk
	X, Y, Z int
}

func (w *World) BlockAt(x, y, z int) *Block {
	ch := w.GetOrGenerateChunk(x>>4, z>>4)
	return &Block{
		World: w,
		Chunk: ch,
		X:     x,
		Y:     y,
		Z:     z,
	}
}

func (b *Block) GetState() uint16 {
	return b.Chunk.GetBlockState(b.X&0xf, b.Y, b.Z&0xf)
}

func (b *Block) SetState(s uint16) {
	b.Chunk.SetBlockState(b.X&0xf, b.Y, b.Z&0xf, s)
}
