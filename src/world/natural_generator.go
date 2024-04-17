package world

import (
	"github.com/aquilax/go-perlin"
	"github.com/cnkei/gospline"
	"gomc/src/data"
	"gomc/src/util"
	"sync"
)

type NaturalGenerator struct {
}

var (
	contSpline   = gospline.NewMonotoneSpline([]float64{-1, -0.3, -0.25, 0.25, 0.3, 0.7, 0.75, 1}, []float64{-0.5, -0.125, 0, 0.25, 0.5, 0.6, 1, 1})
	squashSpline = gospline.NewMonotoneSpline([]float64{-1, 0.25, 0.3, 1}, []float64{3, 6, 30, 40})
)

func (g *NaturalGenerator) Generate(w *World, x, z int) *Chunk {
	blocks := make([]uint16, w.Height<<8)
	p := perlin.NewPerlin(2, 2, 1, w.Seed)
	p2 := perlin.NewPerlin(2, 2, 3, w.Seed)

	wg := sync.WaitGroup{}
	wg.Add(256)
	for bx := 0; bx < 16; bx++ {
		for bz := 0; bz < 16; bz++ {
			bx, bz := bx, bz
			go func() {
				defer wg.Done()

				noise := p.Noise2D(float64((x<<4)+bx)/250.0, float64((z<<4)+bz)/250.0)
				noise += p.Noise2D(float64((x<<4)+bx)/100.0, float64((z<<4)+bz)/100.0) * 0.25
				noise /= 1.25
				cont := contSpline.At(noise)
				squash := squashSpline.At(noise)
				off := 68 + cont*128

				highest := 0
				for y := w.MinY; y < w.MaxY; y++ {
					thresh := util.Map(float64(y)-off, -squash, +squash, -1, 1)
					if thresh >= 0.99 {
						break
					}
					if thresh <= -0.99 {
						blocks[((y-w.MinY)<<8)|(bz<<4)|bx] = data.Stone{}.Id()
						highest = y
						continue
					}
					noise := p2.Noise3D(float64((x<<4)+bx)/50.0+3e7, float64(y)/50.0, float64((z<<4)+bz)/50.0+3e7)
					noise += p2.Noise3D(float64((x<<4)+bx)/25.0+3e7, float64(y)/25.0, float64((z<<4)+bz)/25.0+3e7) * 0.5
					noise /= 1.5
					if noise > thresh {
						blocks[((y-w.MinY)<<8)|(bz<<4)|bx] = data.Stone{}.Id()
						highest = y
					}
				}

				var (
					surfaceBlock uint16
					crustBlock   uint16
				)
				toff := int(p.Noise2D(float64((x<<4)+bx)/10.0, float64((z<<4)+bz)/10.0) * 7.5)

				if highest > 146+toff {
					surfaceBlock = data.SnowBlock{}.Id()
					crustBlock = data.SnowBlock{}.Id()
				} else if highest > 126+toff {
					surfaceBlock = data.GrassBlock{Snowy: true}.Id()
					crustBlock = data.Dirt{}.Id()
				} else if highest > 46+toff {
					if noise < -0.2 {
						surfaceBlock = data.Sand{}.Id()
						crustBlock = data.Sand{}.Id()
					} else {
						surfaceBlock = data.GrassBlock{}.Id()
						crustBlock = data.Dirt{}.Id()
					}
				} else {
					surfaceBlock = data.Gravel{}.Id()
					crustBlock = data.Gravel{}.Id()
				}

				for y := highest; y >= max(w.MinY, highest-4); y-- {
					if blocks[((y-w.MinY)<<8)|(bz<<4)|bx] == 0 {
						continue
					}
					if y == highest {
						blocks[((y-w.MinY)<<8)|(bz<<4)|bx] = surfaceBlock
					} else {
						blocks[((y-w.MinY)<<8)|(bz<<4)|bx] = crustBlock
					}
				}

				for y := w.MinY; y < 64; y++ {
					if blocks[((y-w.MinY)<<8)|(bz<<4)|bx] == 0 {
						blocks[((y-w.MinY)<<8)|(bz<<4)|bx] = data.Water{}.Id()
					}
				}
			}()
		}
	}
	wg.Wait()
	c := NewChunk(w, x, z, blocks)
	c.CalculateSkyLight()
	return c
}
