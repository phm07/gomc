package world

import (
	"testing"
)

func BenchmarkNaturalGenerator(b *testing.B) {
	g := &NaturalGenerator{}
	w := NewWorld(384, -64, 0, g)
	for i := 0; i < b.N; i++ {
		g.Generate(w, i, 0)
	}
}
