package main

import (
	"fmt"
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

func TestComplexitySum(t *testing.T) {
	D := []int{0, 1, 2, 4, 5}
	cs := newComplexitySum(1, D)
	prevlb := 0
	for n := 1; n <= 75; n++ {
		c := 0
		for _, d := range D {
			c += complexity(n + d)
		}

		lb := cs.lowerBoundA(n)
		if lb > c {
			t.Fatalf("lower bound %d exceeds complexity %d", lb, c)
		}
		if lb < prevlb {
			t.Fatalf("lower bound %d decreased below previous %d", lb, prevlb)
		}
		prevlb = lb

		// fmt.Println(n, c, lb)
	}
	fmt.Println(cs.m)
}

func BenchmarkComplexitySum(b *testing.B) {
	D := []int{0, 1, 2, 4, 5}
	for i := 0; i < b.N; i++ {
		cs := newComplexitySum(1, D)
		for n := 1; n <= 75; n++ {
			cs.lowerBoundA(n)
		}
	}
}
