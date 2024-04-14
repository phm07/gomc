package world

import (
	"bytes"
	"gomc/src/nbt"
	"gomc/src/protocol/types"
	"gomc/src/util"
)

type Chunk struct {
	X, Z     int
	Height   int
	Data     []uint16
	SkyLight []byte
}

func NewChunk(height, x, z int) *Chunk {
	return &Chunk{
		Height:   height,
		X:        x,
		Z:        z,
		Data:     make([]uint16, height<<8),
		SkyLight: make([]byte, height<<7+(1<<12)),
	}
}

func (c *Chunk) SetSkyLight(x, y, z int, light byte) {
	idx := ((y << 8) + (z << 4) + x) >> 1
	if (x & 1) == 0 {
		c.SkyLight[idx] = (c.SkyLight[idx] & 0xf0) | (light & 0x0f)
	} else {
		c.SkyLight[idx] = (c.SkyLight[idx] & 0x0f) | ((light & 0x0f) << 4)
	}
}

func (c *Chunk) GetSkyLight(x, y, z int) byte {
	idx := ((y << 8) + (z << 4) + x) >> 1
	if (x & 1) == 0 {
		return c.SkyLight[idx] & 0x0f
	}
	return c.SkyLight[idx] >> 4
}

func (c *Chunk) CalculateSkyLight() {
	obstructed := make([]bool, 256)
	notObstructed := 256
	for i := len(c.Data) - 1; i >= 0 && notObstructed > 0; i-- {
		xz := i & 0xff
		if obstructed[xz] {
			continue
		}
		light := byte(0xf)
		idx := i>>1 + 2048
		if (xz & 1) == 0 {
			c.SkyLight[idx] = (c.SkyLight[idx] & 0xf0) | (light & 0x0f)
		} else {
			c.SkyLight[idx] = (c.SkyLight[idx] & 0x0f) | ((light & 0x0f) << 4)
		}
		if c.Data[i] != 0 {
			obstructed[xz] = true
			notObstructed--
		}
	}
}

func (c *Chunk) MarshalSkyLight() []byte {
	var buf bytes.Buffer
	buf.Write(types.VarInt(len(c.SkyLight) >> 11).Marshal())
	for i := 0; i < len(c.SkyLight)>>11; i++ {
		buf.Write(types.VarInt(2048).Marshal())
		buf.Write(c.SkyLight[i<<11 : (i+1)<<11])
	}
	return buf.Bytes()
}

func (c *Chunk) SetBlockState(x, y, z int, block uint16) {
	c.Data[(y<<8)+(z<<4)+x] = block
}

func (c *Chunk) GetBlockState(x, y, z int) uint16 {
	return c.Data[(y<<8)+(z<<4)+x]
}

func (c *Chunk) Marshal() []byte {
	var buf bytes.Buffer
	for i := 0; i < (c.Height >> 4); i++ {
		section := c.Data[(i << 12):((i + 1) << 12)]
		nonAirBlocks, blockStates := packSection(section, 4, 8)
		biomes := &PalettedContainer{
			Palette: &PaletteSingleValued{
				Value: 39, // plains
			},
		}
		buf.Write(util.Int16ToBytes(int16(nonAirBlocks)))
		buf.Write(blockStates.Marshal())
		buf.Write(biomes.Marshal())
	}
	return buf.Bytes()
}

type ChunkSection struct {
	BlockCount  types.UShort
	BlockStates *PalettedContainer
}

type PalettedContainer struct {
	BitsPerEntry types.Byte
	Palette      Palette
	Data         []uint64
}

func (p *PalettedContainer) Marshal() []byte {
	var buf bytes.Buffer
	buf.Write(p.BitsPerEntry.Marshal())
	buf.Write(p.Palette.Marshal())
	buf.Write(types.VarInt(len(p.Data)).Marshal())
	for _, v := range p.Data {
		buf.Write(util.Uint64ToBytes(v))
	}
	return buf.Bytes()
}

func (c *Chunk) HeightMap() nbt.Tag {
	heightMap := make([]uint16, 256)
	for i := 0; i < (len(c.Data) >> 8); i++ {
		for j := 0; j < 256; j++ {
			if c.Data[(i<<8)+j] != 0 {
				heightMap[j] = uint16(i + 1)
			}
		}
	}
	bpe := util.Log2(c.Height + 1)
	packed := pack(heightMap, bpe)
	tags := make([]*nbt.LongTag, len(packed))
	for i, v := range packed {
		tags[i] = &nbt.LongTag{
			Data: int64(v),
		}
	}
	return &nbt.CompoundTag{
		Data: []nbt.Tag{
			&nbt.ListTag[*nbt.LongTag]{
				Name: "MOTION_BLOCKING",
				Data: tags,
			},
			&nbt.ListTag[*nbt.LongTag]{
				Name: "WORLD_SURFACE",
				Data: tags,
			},
		},
	}
}

func (c *Chunk) GetHeightAt(x int, z int) int {
	for i := c.Height - 1; i >= 0; i-- {
		if c.GetBlockState(x, i, z) != 0 {
			return i
		}
	}
	return 0
}

func packSection(data []uint16, bpeMin, bpeThresh int) (int, *PalettedContainer) {
	count := make(map[uint16]int)
	for _, v := range data {
		count[v]++
	}
	nonAirBlocks := 0
	for k, v := range count {
		if k != 0 {
			nonAirBlocks += v
		}
	}
	if len(count) == 1 {
		return nonAirBlocks, &PalettedContainer{
			Palette: &PaletteSingleValued{
				Value: types.VarInt(data[0]),
			},
		}
	} else if len(count) <= (1 << bpeThresh) {
		palette := make([]types.VarInt, len(count))
		lookup := make(map[uint16]uint16)
		nextId := types.VarInt(0)
		for k := range count {
			palette[nextId] = types.VarInt(k)
			lookup[k] = uint16(nextId)
			nextId++
		}
		bpe := max(util.Log2(len(count)), bpeMin)
		toPack := make([]uint16, len(data))
		for i, v := range data {
			toPack[i] = lookup[v]
		}
		packed := pack(toPack, bpe)
		return nonAirBlocks, &PalettedContainer{
			BitsPerEntry: types.Byte(bpe),
			Palette: &PaletteIndirect{
				Length:  types.VarInt(len(palette)),
				Palette: palette,
			},
			Data: packed,
		}
	} else {
		packed := pack(data, 15)
		return nonAirBlocks, &PalettedContainer{
			BitsPerEntry: 15,
			Palette:      &PaletteDirect{},
			Data:         packed,
		}
	}
}

func pack(data []uint16, bpe int) []uint64 {
	res := make([]uint64, len(data)/(64/bpe))
	idx, shift := 0, 0
	for _, v := range data {
		if shift > 64-bpe {
			shift = 0
			idx++
		}
		res[idx] |= uint64(v) << shift
		shift += bpe
	}
	return res
}
