package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func BenchmarkMelody(b *testing.B) {
	rand.Seed(time.Now().UnixNano())

	m := newMelody(1, 8)
	for i := 0; i < b.N; i++ {
		m.time += .25
		m.next(allRats)
	}
	fmt.Println(m.history[len(m.history)-1].n)
}
