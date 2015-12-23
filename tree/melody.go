package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
)

type melody struct {
	center        float64
	coherency     float64
	coherencyTime float64

	time    float64
	history []*note
}

type note struct {
	time interval
	abs  float64
	rel  int
}

func newMelody(center, coherencyTime float64) melody {
	return melody{
		center:        center,
		coherency:     math.Pow(.01, -1/coherencyTime),
		coherencyTime: coherencyTime,
		history:       []*note{{interval{}, center, 1}},
	}
}

func (m *melody) next(duration float64) (*note, ratio) {
	return m.nextAfter(duration, m.history[len(m.history)-1], allRats)
}

func (m *melody) nextAfter(duration float64, prev *note, rats []ratio) (*note, ratio) {
	if prev.time.max < m.time-m.coherencyTime {
		fmt.Printf("melody: %.2f < %.2f\n", prev.time.max, m.time-m.coherencyTime)
	}

	time := interval{m.time, m.time + duration}

	cSum, ampSum := m.historyComplexity(time)
	sum := 0.0
	sums := make([]float64, len(rats))
	for i, r := range rats {
		p := math.Log2(prev.abs * r.float() / m.center)
		sum += math.Exp2(-p*p/2) * math.Exp2(-m.complexity(time, prev, cSum, ampSum, r))
		sums[i] = sum
	}
	i := sort.SearchFloat64s(sums, sum*rand.Float64())
	return m.appendHistory(time, prev, rats[i]), rats[i]
}

func (m *melody) historyComplexity(time interval) (cSum, ampSum float64) {
	for i, n1 := range m.history {
		a1 := math.Pow(m.coherency, time.overlap(n1.time))
		for _, n2 := range m.history[:i] {
			a2 := math.Pow(m.coherency, time.overlap(n2.time))
			cSum += a1 * a2 * float64(complexity(n1.rel, n2.rel))
		}
		ampSum += a1
	}
	return
}

func (m *melody) complexity(time interval, prev *note, cSum, ampSum float64, r ratio) float64 {
	const a1 = 1
	n1n := r.a * prev.rel
	for _, n2 := range m.history {
		a2 := math.Pow(m.coherency, time.overlap(n2.time))
		cSum += a1 * a2 * float64(complexity(n1n, n2.rel*r.b))
	}
	return cSum / (ampSum + a1)
}

func complexity(a, b int) int {
	c := 0
	for d := 2; a != b; {
		d1 := a%d == 0
		d2 := b%d == 0
		if d1 != d2 {
			c += d - 1
		}
		if d1 {
			a /= d
		}
		if d2 {
			b /= d
		}
		if !(d1 || d2) {
			d++
		}
	}
	return c
}

func (m *melody) appendHistory(time interval, prev *note, r ratio) *note {
	n := &note{time, prev.abs * r.float(), prev.rel * r.a}
	for i := range m.history {
		m.history[i].rel *= r.b
	}
	m.history = append(m.history, n)

	newhist := make([]*note, 0, len(m.history))
	for _, n := range m.history {
		if n.time.max >= m.time-m.coherencyTime {
			newhist = append(newhist, n)
		}
	}
	m.history = newhist

	d := m.history[0].rel
	for _, n := range m.history[1:] {
		d = gcd(d, n.rel)
	}
	for i := range m.history {
		m.history[i].rel /= d
	}

	return n
}

func gcd(a, b int) int {
	for a > 0 {
		if a > b {
			a, b = b, a
		}
		b -= a
	}
	return b
}

type interval struct {
	min, max float64
}

// overlap reports the amount of overlap between x and y, or, if there is no overlap, their negative separation.
func (x interval) overlap(y interval) float64 {
	return math.Min(
		math.Min(
			x.max-x.min,
			y.max-y.min,
		),
		math.Min(
			x.max-y.min,
			y.max-x.min,
		),
	)
}

type ratio struct {
	a, b int
}

func (r ratio) float() float64 { return float64(r.a) / float64(r.b) }

var allRats = func() (rats []ratio) {
	pow := func(a, x int) int {
		y := 1
		for x > 0 {
			y *= a
			x--
		}
		return y
	}
	mul := func(n, d, a, x int) (int, int) {
		if x > 0 {
			return n * pow(a, x), d
		}
		return n, d * pow(a, -x)
	}
	for _, two := range []int{-3, -2, -1, 0, 1, 2, 3} {
		for _, three := range []int{-2, -1, 0, 1, 2} {
			for _, five := range []int{-1, 0, 1} {
				for _, seven := range []int{-1, 0, 1} {
					n, d := 1, 1
					n, d = mul(n, d, 2, two)
					n, d = mul(n, d, 3, three)
					n, d = mul(n, d, 5, five)
					n, d = mul(n, d, 7, seven)
					if complexity(n, d) < 12 {
						rats = append(rats, ratio{n, d})
					}
				}
			}
		}
	}
	return
}()
