package world

import (
	"testing"
)

func BenchmarkNaturalGenerator(b *testing.B) {
	g := &NaturalGenerator{}
	for i := 0; i < b.N; i++ {
		g.Generate(0, 384, i, 0)
	}
}
