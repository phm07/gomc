package world

type World struct {
	Height     int
	MinY, MaxY int
	Chunks     map[int64]*Chunk
	Generator  Generator
	Seed       int64
}

func NewWorld(height int, minY int, seed int64, generator Generator) *World {
	return &World{
		Height:    height,
		MinY:      minY,
		MaxY:      height + minY - 1,
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
		ch = w.Generator.Generate(w, x, z)
		w.Chunks[int64(x)<<32|(int64(z)&0xffffffff)] = ch
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
