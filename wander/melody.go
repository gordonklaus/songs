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
	rhythm        bool
	bias          float64
	center        float64
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
	complexity := .5 // 0..1
	return melody{
		bias:          math.Log2(complexity),
		center:        center,
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
	notes := m.history
	if m.rhythm {
		notes = make([]note, 0, len(m.history)*(1+len(m.history))/2)
		for i := range m.history {
			n := 0
			for _, n2 := range m.history[i:] {
				n += n2.n
				notes = append(notes, note{n: n})
			}
		}
	}

	cSum := m.complexitySum(notes)

	rats := rats
	if m.rhythm {
		// rats = append(rats, ratio{0, 1})
	}
	sum := 0.0
	sums := make([]float64, len(rats))
	for i, r := range rats {
		f := r.float()
		if f == 0 {
			f = 1
		}
		p := math.Log2(m.last * f / m.center)
		sum += math.Exp2(-p*p/2) * math.Exp2(m.bias*m.complexity(notes, cSum, r))
		sums[i] = sum
	}
	i := 0
	x := sum * rand.Float64()
	for i = range sums {
		if x < sums[i] {
			break
		}
	}
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

	if rats[i].a == 0 {
		return 0
	} else {
		m.last *= rats[i].float()
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

func (m *melody) complexitySum(notes []note) int {
	c := 0
	for i, n1 := range notes {
		for _, n2 := range notes[:i] {
			c += complexity(n1.n, n2.n)
		}
	}
	return c
}

func (m *melody) complexity(notes []note, cSum int, r ratio) float64 {
	last := 1
	for i := len(m.history)-1; i >= 0; i-- {
		if m.history[i].n != 0 {
			last = m.history[i].n
			break
		}
	}
	n := r.a * last
	notes2 := []note{{n: n}}
	if m.rhythm {
		notes2 = make([]note, 1 + len(m.history))
		notes2[len(m.history)] = note{n: n}
		for i := len(m.history) - 1; i >= 0; i-- {
			n += m.history[i].n * r.b
			notes2[i] = note{n: n}
		}
	}
	for i, n1 := range notes2 {
		for _, n2 := range notes2[:i] {
			cSum += complexity(n1.n, n2.n)
		}
	}
	for _, n1 := range notes {
		for _, n2 := range notes2 {
			cSum += complexity(n1.n*r.b, n2.n)
		}
	}
	return float64(cSum) / float64(len(notes) + len(notes2))
}

func complexity(a, b int) int {
	if a == 0 || b == 0 {
		return 0
	}
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
	last := 1
	for i := len(m.history)-1; i >= 0; i-- {
		if m.history[i].n != 0 {
			last = m.history[i].n
			break
		}
	}
	r.a *= last
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
