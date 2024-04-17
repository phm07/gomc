package world

import (
	"bytes"
	"gomc/src/protocol/types"
	"gomc/src/util"
	"sync"
)

type PalettedContainer struct {
	Size      int
	BpeMin    int
	BpeThresh int
	Count     map[uint16]int
	Bpe       int
	Palette   Palette
	Data      []uint64
}

func (p *PalettedContainer) GetDataAt(idx int) uint16 {
	switch pal := p.Palette.(type) {
	case *PaletteSingleValued:
		return uint16(pal.Value)
	case *PaletteDirect:
		return p.getPackedData(idx)
	case *PaletteIndirect:
		return uint16(pal.Palette[p.getPackedData(idx)])
	default:
		panic("unknown palette")
	}
}

func (p *PalettedContainer) SetDataAt(idx int, v uint16) {
	switch pal := p.Palette.(type) {
	case *PaletteSingleValued:
		prev := uint16(pal.Value)
		if prev == v {
			return
		}
		p.Count[prev]--
		p.Count[v]++
		data := p.Uncompress()
		data[idx] = v
		p.Palette, p.Bpe, p.Data = packPaletted(data, p.BpeMin)
	case *PaletteDirect:
		prev := p.getPackedData(idx)
		if prev == v {
			return
		}
		p.Count[prev]--
		if p.Count[prev] == 0 {
			delete(p.Count, prev)
		}
		if len(p.Count) <= (1 << p.BpeThresh) {
			data := p.Uncompress()
			data[idx] = v
			p.Palette, p.Bpe, p.Data = packPaletted(data, p.BpeMin)
		} else {
			p.Count[v]++
			p.setPackedData(idx, v)
		}
	case *PaletteIndirect:
		prevIdx := p.getPackedData(idx)
		prev := uint16(pal.Palette[prevIdx])
		if prev == v {
			return
		}
		p.Count[prev]--
		if p.Count[prev] == 0 {
			delete(p.Count, prev)
		}
		p.Count[v]++
		if len(p.Count) <= 1 {
			v := uint16(0)
			for k := range p.Count {
				v = k
			}
			p.Palette, p.Bpe, p.Data = &PaletteSingleValued{Value: types.VarInt(v)}, 0, nil
			return
		}
		if len(p.Count) > (1 << p.BpeThresh) {
			data := p.Uncompress()
			data[idx] = v
			p.Palette, p.Bpe, p.Data = &PaletteDirect{}, 15, pack(data, 15)
			return
		}
		newBpe := max(getBpe(len(p.Count)), p.BpeMin)
		if p.Bpe != newBpe {
			data := unpack(p.Data, p.Bpe, p.Size)
			if p.Count[prev] == 0 {
				pal.Palette = append(pal.Palette[:prevIdx], pal.Palette[(prevIdx+1):]...)
				for i, v := range data {
					if v > prevIdx {
						data[i]--
					}
				}
			}
			if p.Count[v] == 1 {
				pal.Palette = append(pal.Palette, types.VarInt(v))
				data[idx] = uint16(len(pal.Palette) - 1)
			}
			p.Data, p.Bpe = pack(data, newBpe), newBpe
			return
		}
		if p.Count[prev] > 0 {
			if vIdx := pal.IndexOf(types.VarInt(v)); vIdx > -1 {
				p.setPackedData(idx, uint16(vIdx))
			} else {
				pal.Palette = append(pal.Palette, types.VarInt(v))
				p.setPackedData(idx, uint16(len(pal.Palette)-1))
			}
		} else if p.Count[v] == 1 {
			pal.Palette[prevIdx] = types.VarInt(v)
		} else {
			pal.Palette = append(pal.Palette[:prevIdx], pal.Palette[(prevIdx+1):]...)
			data := unpack(p.Data, p.Bpe, p.Size)
			for i, v := range data {
				if v > prevIdx {
					data[i]--
				}
			}
			data[idx] = uint16(pal.IndexOf(types.VarInt(v)))
			p.Data = pack(data, p.Bpe)
		}
	}
}

func (p *PalettedContainer) getPackedData(idx int) uint16 {
	dpl := 64 / p.Bpe
	return uint16((p.Data[idx/dpl] >> (p.Bpe * (idx % dpl))) & ((1 << p.Bpe) - 1))
}

func (p *PalettedContainer) setPackedData(idx int, v uint16) {
	bpe := int(p.Bpe)
	dpl := 64 / bpe
	i, m, off := idx/dpl, uint16((1<<bpe)-1), (idx%dpl)*bpe
	p.Data[i] &^= uint64(m) << off
	p.Data[i] |= uint64(v&m) << off
}

func (p *PalettedContainer) Marshal() []byte {
	var buf bytes.Buffer
	buf.Write(types.Byte(p.Bpe).Marshal())
	buf.Write(p.Palette.Marshal())
	buf.Write(types.VarInt(len(p.Data)).Marshal())
	for _, v := range p.Data {
		buf.Write(util.Uint64ToBytes(v))
	}
	return buf.Bytes()
}

func (p *PalettedContainer) Uncompress() []uint16 {
	switch pal := p.Palette.(type) {
	case *PaletteSingleValued:
		data := make([]uint16, p.Size)
		if v := uint16(pal.Value); v != 0 {
			for i := 0; i < len(data); i++ {
				data[i] = v
			}
		}
		return data
	case *PaletteDirect:
		return unpack(p.Data, 15, p.Size)
	case *PaletteIndirect:
		lookup := make(map[uint16]uint16)
		for k, v := range pal.Palette {
			lookup[uint16(k)] = uint16(v)
		}
		data := unpack(p.Data, p.Bpe, p.Size)
		for i := 0; i < len(data); i++ {
			v := data[i]
			data[i] = lookup[v]
			for ; i < len(data)-1 && data[i+1] == v; i++ {
				data[i+1] = data[i]
			}
		}
		return data
	default:
		panic("unknown palette")
	}
}

func PalettedContainerFromBytes(data []uint16, bpeMin, bpeThresh int) *PalettedContainer {
	count := make(map[uint16]int)
	for _, v := range data {
		count[v]++
	}
	p := &PalettedContainer{
		Size:      len(data),
		BpeMin:    bpeMin,
		BpeThresh: bpeThresh,
		Count:     count,
	}
	if len(count) == 1 {
		p.Palette = &PaletteSingleValued{Value: types.VarInt(data[0])}
	} else if len(count) <= (1 << bpeThresh) {
		p.Palette, p.Bpe, p.Data = packPaletted(data, bpeMin)
	} else {
		p.Palette, p.Bpe, p.Data = &PaletteDirect{}, 15, pack(data, 15)
	}
	return p
}

func PalettedContainersFromBytes(data []uint16, bpeMin, bpeThresh int) []*PalettedContainer {
	sections := make([]*PalettedContainer, len(data)>>12)
	wg := sync.WaitGroup{}
	wg.Add(len(sections))
	for i := 0; i < len(sections); i++ {
		go func() {
			defer wg.Done()
			sections[i] = PalettedContainerFromBytes(data[(i<<12):((i+1)<<12)], bpeMin, bpeThresh)
		}()
	}
	wg.Wait()
	return sections
}

func pack(data []uint16, bpe int) []uint64 {
	dpl := 64 / bpe
	res := make([]uint64, (len(data)+dpl-1)/dpl)
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

func packPaletted(data []uint16, bpeMin int) (Palette, int, []uint64) {
	var palette []types.VarInt
	lookup := make(map[uint16]uint16)
	toPack := make([]uint16, len(data))
	var ok bool
	for i := 0; i < len(data); i++ {
		v := data[i]
		if toPack[i], ok = lookup[v]; !ok {
			palette = append(palette, types.VarInt(v))
			k := uint16(len(palette) - 1)
			lookup[v] = k
			toPack[i] = k
		}
		for ; i < len(data)-1 && data[i+1] == v; i++ {
			toPack[i+1] = toPack[i]
		}
	}
	bpe := max(getBpe(len(palette)), bpeMin)
	packed := pack(toPack, bpe)
	return &PaletteIndirect{
		Length:  types.VarInt(len(palette)),
		Palette: palette,
	}, bpe, packed
}

func unpack(data []uint64, bpe int, n int) []uint16 {
	dpl := 64 / bpe
	res := make([]uint16, 0, n)
	mask := uint16((1 << bpe) - 1)
	for _, v := range data {
		for j := 0; j < dpl && n > 0; j++ {
			res = append(res, uint16(v>>(j*bpe))&mask)
			n--
		}
	}
	return res
}

func getBpe(n int) int {
	return util.Log2(n-1) + 1
}
