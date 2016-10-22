package main

import (
	"math/rand"
	"testing"
	"time"
)

func BenchmarkMelody(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	m := NewMelody()
	for i := 0; i < b.N; i++ {
		m.Next()
	}
}
