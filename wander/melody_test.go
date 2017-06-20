package main

import (
	"math/rand"
	"testing"
	"time"
)

func BenchmarkMelody(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	m := NewMelody()
	for i := 0; i < 10; i++ {
		m.Next()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Next()
	}
}

func TestMelody(t *testing.T) {
	m := NewMelody()
	for i := 0; ; i++ {
		m.Next()
	}
}
