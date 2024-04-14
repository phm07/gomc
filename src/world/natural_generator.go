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

func (g *NaturalGenerator) Generate(height, x, z int) *Chunk {
	c := NewChunk(height, x, z)
	p := perlin.NewPerlin(2, 2, 1, 0)
	p2 := perlin.NewPerlin(2, 2, 3, 0)

	wg := sync.WaitGroup{}
	wg.Add(256)
	for bx := 0; bx < 16; bx++ {
		for bz := 0; bz < 16; bz++ {
			bx, bz := bx, bz
			go func() {
				defer wg.Done()

				noise := p.Noise2D(float64((x<<4)+bx)/250.0, float64((z<<4)+bz)/250.0)
				noise += p.Noise2D(float64((x<<4)+bx)/100.0, float64((z<<4)+bz)/100.0) * 0.25
				cont := contSpline.At(noise)
				squash := squashSpline.At(noise)
				off := 132 + cont*128

				highest := 0
				for y := 0; y < c.Height; y++ {
					thresh := util.Map(float64(y)-off, -squash, +squash, -1, 1)
					if thresh >= 0.99 {
						break
					}
					if thresh <= -0.99 {
						c.SetBlockState(bx, y, bz, data.Stone{}.Id())
						highest = y
						continue
					}
					noise := p2.Noise3D(float64((x<<4)+bx)/50.0+3e7, float64(y)/50.0, float64((z<<4)+bz)/50.0+3e7)
					noise += p2.Noise3D(float64((x<<4)+bx)/25.0+3e7, float64(y)/25.0, float64((z<<4)+bz)/25.0+3e7) * 0.5
					noise = max(min(1, noise), -1)
					if noise > thresh {
						c.SetBlockState(bx, y, bz, data.Stone{}.Id())
						highest = y
					}
				}

				var (
					surfaceBlock uint16
					crustBlock   uint16
				)
				toff := int(p.Noise2D(float64((x<<4)+bx)/10.0, float64((z<<4)+bz)/10.0) * 7.5)

				if highest > 210+toff {
					surfaceBlock = data.SnowBlock{}.Id()
					crustBlock = data.SnowBlock{}.Id()
				} else if highest > 190+toff {
					surfaceBlock = data.GrassBlock{Snowy: true}.Id()
					crustBlock = data.Dirt{}.Id()
				} else if highest > 110+toff {
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

				for y := highest; y >= max(0, highest-4); y-- {
					if c.GetBlockState(bx, y, bz) == 0 {
						continue
					}
					if y == highest {
						c.SetBlockState(bx, y, bz, surfaceBlock)
					} else {
						c.SetBlockState(bx, y, bz, crustBlock)
					}
				}

				for y := 0; y < 128; y++ {
					if c.GetBlockState(bx, y, bz) == 0 {
						c.SetBlockState(bx, y, bz, data.Water{}.Id())
					}
				}
			}()
		}
	}
	wg.Wait()
	c.CalculateSkyLight()
	return c
}
