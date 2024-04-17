package world

import (
	"bytes"
	"gomc/src/nbt"
	"gomc/src/protocol/packet"
	"gomc/src/protocol/types"
	"gomc/src/util"
)

type Chunk struct {
	X, Z     int
	World    *World
	Sections []*PalettedContainer
	SkyLight []byte
}

func NewChunk(w *World, x, z int, data []uint16) *Chunk {
	return &Chunk{
		World:    w,
		X:        x,
		Z:        z,
		Sections: PalettedContainersFromBytes(data, 4, 8),
		SkyLight: make([]byte, w.Height<<7+(1<<12)),
	}
}

func (c *Chunk) Packet() packet.SerializablePacket {
	mask := types.NewBitSet(26)
	mask2 := types.NewBitSet(26)
	for i := 0; i < 26; i++ {
		mask2.SetBit(i, true)
	}
	return &packet.ClientboundPlayChunkData{
		ChunkX:               types.Int(c.X),
		ChunkZ:               types.Int(c.Z),
		Heightmaps:           c.HeightMap().Marshal(),
		Data:                 c.Marshal(),
		NumBlockEntities:     0,
		SkyLightMask:         mask2,
		BlockLightMask:       mask,
		EmptySkyLightMask:    mask,
		EmptyBlockLightMask:  mask2,
		SkyLight:             c.MarshalSkyLight(),
		BlockLightArrayCount: types.VarInt(0).Marshal(),
	}
}

func (c *Chunk) SetSkyLight(x, y, z int, light byte) {
	y -= c.World.MinY
	idx := ((y << 8) + (z << 4) + x) >> 1
	if (x & 1) == 0 {
		c.SkyLight[idx] = (c.SkyLight[idx] & 0xf0) | (light & 0x0f)
	} else {
		c.SkyLight[idx] = (c.SkyLight[idx] & 0x0f) | ((light & 0x0f) << 4)
	}
}

func (c *Chunk) GetSkyLight(x, y, z int) byte {
	y -= c.World.MinY
	idx := ((y << 8) + (z << 4) + x) >> 1
	if (x & 1) == 0 {
		return c.SkyLight[idx] & 0x0f
	}
	return c.SkyLight[idx] >> 4
}

func (c *Chunk) CalculateSkyLight() {
	obstructed := make([]bool, 256)
	notObstructed := 256
	data := c.blockData()
	for i := len(data) - 1; i >= 0 && notObstructed > 0; i-- {
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
		if data[i] != 0 {
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
	y -= c.World.MinY
	c.Sections[y>>12].SetDataAt(((y&0xf)<<8)+(z<<4)+x, block)
}

func (c *Chunk) GetBlockState(x, y, z int) uint16 {
	y -= c.World.MinY
	return c.Sections[y>>12].GetDataAt(((y & 0xf) << 8) + (z << 4) + x)
}

func (c *Chunk) Marshal() []byte {
	var buf bytes.Buffer
	for _, s := range c.Sections {
		var nonAirBlocks int16
		for k, v := range s.Count {
			if k > 0 {
				nonAirBlocks += int16(v)
			}
		}
		biomes := &PalettedContainer{
			Palette: &PaletteSingleValued{
				Value: 39, // plains
			},
		}
		buf.Write(util.Int16ToBytes(nonAirBlocks))
		buf.Write(s.Marshal())
		buf.Write(biomes.Marshal())
	}
	return buf.Bytes()
}

type ChunkSection struct {
	BlockCount  types.UShort
	BlockStates *PalettedContainer
}

func (c *Chunk) HeightMap() nbt.Tag {
	heightMap := make([]uint16, 256)
	data := c.blockData()
	for i := 0; i < (len(data) >> 8); i++ {
		for j := 0; j < 256; j++ {
			if data[(i<<8)+j] != 0 {
				heightMap[j] = uint16(i + 1 + c.World.MinY)
			}
		}
	}
	bpe := util.Log2(c.World.Height + 1)
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
	for y := c.World.MaxY; y >= c.World.MinY; y-- {
		if c.GetBlockState(x, y, z) != 0 {
			return y
		}
	}
	return c.World.MinY
}

func (c *Chunk) blockData() []uint16 {
	var data []uint16
	for _, s := range c.Sections {
		data = append(data, s.Uncompress()...)
	}
	return data
}
