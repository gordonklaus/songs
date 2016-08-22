package main

import (
	"math"
	"math/rand"
)

type ratio struct {
	a, b int
}

func (r ratio) float() float64 { return float64(r.a) / float64(r.b) }

type melody struct {
	bias          float64
	rhythm        bool
	center        float64
	coherency     float64
	coherencyTime float64

	last    float64
	time    float64
	history []note
}

type note struct {
	t float64
	n int
}

func newMelody(center, coherencyTime float64) melody {
	complexity := .6 // 0..1
	return melody{
		bias:          math.Log2(complexity),
		center:        center,
		coherency:     math.Pow(.01, 1./coherencyTime),
		coherencyTime: coherencyTime,
		last:          center,
		history:       []note{{0, 1}},
	}
}

func newRhythm(center, coherencyTime float64) melody {
	m := newMelody(center, coherencyTime)
	m.rhythm = true
	return m
}

func (m *melody) next() float64 {
	cSum, ampSum := m.historyComplexity()

	sum := 0.0
	sums := make([]float64, len(rats))
	for i, r := range rats {
		p := math.Log2(m.last * r.float() / m.center)
		sum += math.Exp2(-p*p/2) * math.Exp2(m.bias*m.complexity(cSum, ampSum, r))
		sums[i] = sum
	}
	if m.rhythm {
		sum += math.Exp2(m.bias * cSum / ampSum)
		sums = append(sums, sum)
	}
	i := 0
	x := sum * rand.Float64()
	for i = range sums {
		if x < sums[i] {
			break
		}
	}
	if i == len(rats) {
		return 0
	}
	m.last *= rats[i].float()
	m.history = m.appendHistory(rats[i])

	for i, n := range m.history {
		if m.time-n.t < m.coherencyTime {
			m.history = m.history[i:]
			d := m.history[0].n
			for _, n := range m.history[1:] {
				d = gcd(d, n.n)
			}
			for i := range m.history {
				m.history[i].n /= d
			}
			break
		}
	}

	return m.last
}

var rats []ratio

func init() {
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
}

func (m *melody) historyComplexity() (cSum, ampSum float64) {
	for i, n1 := range m.history {
		a1 := math.Pow(m.coherency, m.time-n1.t)
		for _, n2 := range m.history[:i] {
			a2 := math.Pow(m.coherency, m.time-n2.t)
			cSum += a1 * a2 * float64(complexity(n1.n, n2.n))
		}
		ampSum += a1
	}
	return
}

func (m *melody) complexity(cSum, ampSum float64, r ratio) float64 {
	const a1 = 1
	n1n := r.a * m.history[len(m.history)-1].n
	for _, n2 := range m.history {
		a2 := math.Pow(m.coherency, m.time-n2.t)
		cSum += a1 * a2 * float64(complexity(n1n, n2.n*r.b))
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

func (m *melody) appendHistory(r ratio) []note {
	r.a *= m.history[len(m.history)-1].n
	history := make([]note, len(m.history), len(m.history)+1)
	for i, n := range m.history {
		history[i] = note{n.t, n.n * r.b}
	}
	return append(history, note{m.time, r.a})
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
