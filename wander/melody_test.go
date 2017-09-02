package main

import (
	"fmt"
	"math"
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
	lbi := getLowerBound(1, 0, D)
	prevc := 0
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
		if lb > prevlb && prevlb != prevc {
			t.Fatalf("lower bound increased to %d but previous lower bound %d did not equal previous complexity %c", lb, prevlb, prevc)
		}

		if n >= lbi.n1 {
			lbi.increment()
		}
		lb2 := lbi.value
		if lb2 != lb {
			t.Fatalf("lb2=%d lb1=%d\n%v\n%v\n%d", lb2, lb, lbi.lb.lb.steps, lbi.lb.lb.pending, lbi.lb.lb.m)
		}
		fmt.Println(n, c, lb, lb2)

		prevlb = lb
		prevc = c
	}
	fmt.Println(cs.m, cs.l)
	fmt.Println(cs.lb)
	fmt.Println(lbi.lb.lb.m, lbi.lb.lb.l)
	fmt.Println(lbi.lb.lb.pending)
}

func TestProbabilitySum(t *testing.T) {
	sum := 1.0
	for i := 0; i < 25; i++ {
		p := prime(i)
		x := math.Pow(2, float64(p-1))
		sum *= (x + 1) / (x - 1)
		fmt.Println(i, sum)
	}
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

type complexitySum struct {
	b      int
	D      []int
	dhmean float64
	lb     []lowerBoundStepOld
	m, l   int
}

type lowerBoundStepOld struct {
	n, c int
}

func newComplexitySum(b int, D []int) *complexitySum {
	dprod := 1
	count := 0
	for _, d := range D {
		if d > 0 {
			dprod *= d
			count++
		}
	}
	cs := &complexitySum{
		D:      D,
		dhmean: 1 / math.Pow(float64(dprod), 1/float64(count)),
		m:      1,
		l:      1,
	}
	return cs
}

func (cs *complexitySum) lowerBoundA(n int) int {
	if len(cs.lb) > 1 && n >= cs.lb[1].n {
		cs.lb = cs.lb[1:]
		cs.l = int(math.Ceil(math.Exp2(float64(cs.lb[0].c) / float64(len(cs.D)))))
	}
	// fmt.Printf("n=%d\n", n)
	for ; cs.m <= n || cs.m <= cs.l; cs.m++ {
		// if gcd(cs.m, cs.b) != 1 {
		// 	continue
		// }
		// fmt.Printf("\tm=%d <= l=%d\n", cs.m, cs.l)
		c := 0
		for _, d := range cs.D {
			c += complexity(cs.m + d)
		}

		i := len(cs.lb)
		for ; i > 0; i-- {
			if c > cs.lb[i-1].c {
				break
			}
		}
		if i < len(cs.lb) {
			cs.lb = cs.lb[:i+1]
			cs.lb[i].c = c
		} else {
			cs.lb = append(cs.lb, lowerBoundStepOld{cs.m, c})
			if len(cs.lb) > 1 && n >= cs.lb[1].n {
				cs.lb = cs.lb[1:]
				i = 0
			}
		}
		if i == 0 {
			cs.l = int(math.Ceil(math.Exp2(float64(cs.lb[0].c) / float64(len(cs.D)))))
		}
		// fmt.Printf("\tlb=%v l=%d\n", cs.lb, cs.l)
	}

	// fmt.Println()
	return cs.lb[0].c
}

func (cs *complexitySum) lowerBoundB(n int) int {
	if len(cs.lb) > 1 && n >= cs.lb[1].n {
		cs.lb = cs.lb[1:]
		cs.l = int(math.Ceil(math.Exp2(float64(cs.lb[0].c)/float64(len(cs.D)))) * cs.dhmean)
	}
	// fmt.Printf("n=%d\n", n)
	for ; cs.m <= n || cs.m <= cs.l; cs.m++ {
		// fmt.Printf("\tm=%d <= l=%d\n", cs.m, cs.l)
		c := 0
		for _, d := range cs.D {
			c += complexity(1 + cs.m*d)
		}

		i := len(cs.lb)
		for ; i > 0; i-- {
			if c > cs.lb[i-1].c {
				break
			}
		}
		if i < len(cs.lb) {
			cs.lb = cs.lb[:i+1]
			cs.lb[i].c = c
		} else {
			cs.lb = append(cs.lb, lowerBoundStepOld{cs.m, c})
			if len(cs.lb) > 1 && n >= cs.lb[1].n {
				cs.lb = cs.lb[1:]
				i = 0
			}
		}
		if i == 0 {
			cs.l = int(math.Ceil(math.Exp2(float64(cs.lb[0].c)/float64(len(cs.D)))) * cs.dhmean)
		}
		// fmt.Printf("\tlb=%v l=%d\n", cs.lb, cs.l)
	}

	// fmt.Println()
	return cs.lb[0].c
}
