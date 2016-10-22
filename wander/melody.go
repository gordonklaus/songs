package main

import (
	"math"
	"math/rand"
	"sort"
)

type Melody struct {
	rhythmBias    float64
	frequencyBias float64
	avgDuration   float64
	avgFrequency  float64
	coherencyTime float64

	lastDuration  float64
	lastFrequency float64
	history       []note
}

type note struct {
	t int
	f int
}

func NewMelody() *Melody {
	rhythmComplexity := .85   // 0..1
	frequencyComplexity := .5 // 0..1
	avgDuration := 0.5
	avgFrequency := 256.0
	coherencyTime := 8.0
	history := make([]note, int(coherencyTime/avgDuration))
	for i := range history {
		history[i] = note{i, 1}
	}
	return &Melody{
		rhythmBias:    math.Log2(rhythmComplexity),
		frequencyBias: math.Log2(frequencyComplexity),
		avgDuration:   avgDuration,
		avgFrequency:  avgFrequency,
		coherencyTime: coherencyTime,
		lastDuration:  avgDuration,
		lastFrequency: avgFrequency,
		history:       history,
	}
}

func (m *Melody) Next() (float64, float64) {
	m.appendHistory(m.newDuration(), m.newFrequency())
	return m.lastDuration, m.lastFrequency
}

func (m *Melody) newDuration() ratio {
	harmonies := make([]int, 0, len(m.history)*(len(m.history)-1)/2)
	for i, n1 := range m.history {
		for _, n0 := range m.history[:i] {
			harmonies = append(harmonies, n1.t-n0.t)
		}
	}
	return selectRatio(func(r ratio) float64 {
		return math.Exp(-m.lastDuration*r.float()/m.avgDuration) * math.Exp2(m.rhythmBias*m.durationComplexity(harmonies, r))
	})
}

func (m *Melody) newFrequency() ratio {
	harmonies := make([]int, len(m.history))
	for i, n := range m.history {
		harmonies[i] = n.f
	}
	return selectRatio(func(r ratio) float64 {
		dp := math.Log2(m.lastFrequency * r.float() / m.avgFrequency)
		return math.Exp2(-dp*dp/2) * math.Exp2(m.frequencyBias*m.frequencyComplexity(harmonies, r))
	})
}

func selectRatio(complexity func(ratio) float64) ratio {
	sum := 0.0
	sums := make([]float64, len(ratios))
	for i, r := range ratios {
		sum += complexity(r)
		sums[i] = sum
	}
	return ratios[sort.SearchFloat64s(sums, sum*rand.Float64())]
}

func (m *Melody) durationComplexity(harmonies []int, r ratio) float64 {
	harmoniesA := make([]int, len(m.history))
	last := m.history[len(m.history)-1]
	last2 := m.history[len(m.history)-2]
	d := r.a * (last.t - last2.t)
	t1 := r.b*last.t + d
	for i, n := range m.history {
		t0 := r.b * n.t
		harmoniesA[i] = t1 - t0
	}
	return m.complexity(harmonies, harmoniesA, r.b)
}

func (m *Melody) frequencyComplexity(harmonies []int, r ratio) float64 {
	harmoniesA := []int{r.a * harmonies[len(harmonies)-1]}
	return m.complexity(harmonies, harmoniesA, r.b)
}

func (m *Melody) complexity(harmonies, harmoniesA []int, b int) float64 {
	c := 0
	for i, h1 := range harmoniesA {
		for _, h2 := range harmoniesA[:i] {
			c += complexity(h1, h2)
		}
	}
	for _, h1 := range harmonies {
		for _, h2 := range harmoniesA {
			c += complexity(h1*b, h2)
		}
	}
	n := len(harmonies) + len(harmoniesA)
	n = n * (n - 1) / len(m.history)
	return float64(c) / float64(n)
}

func (m *Melody) appendHistory(rd, rf ratio) {
	last := m.history[len(m.history)-1]
	last2 := m.history[len(m.history)-2]
	for i := range m.history {
		m.history[i].t *= rd.b
		m.history[i].f *= rf.b
	}
	d := rd.a * (last.t - last2.t)
	t1 := rd.b*last.t + d
	m.history = append(m.history, note{
		t: t1,
		f: rf.a * last.f,
	})

	m.lastDuration *= rd.float()
	m.lastFrequency *= rf.float()

	for i, n := range m.history[:len(m.history)-2] {
		r := ratio{t1 - n.t, d}
		if m.lastDuration*r.float() < m.coherencyTime {
			m.history = m.history[i:]
			break
		}
	}

	t0 := m.history[0].t
	for i := range m.history {
		m.history[i].t -= t0
	}

	td, fd := 0, 0
	for _, n := range m.history {
		td = gcd(td, n.t)
		fd = gcd(fd, n.f)
	}
	for i := range m.history {
		m.history[i].t /= td
		m.history[i].f /= fd
	}
}

var complexityCache = map[int]int{}

func complexity(a, b int) int {
	if a == b || a == 0 || b == 0 {
		return 0
	}
	d := gcd(a, b)
	d *= d
	n := a * b / d
	if c, ok := complexityCache[n]; ok {
		return c
	}
	c := 0
	for m, d := 1, 2; m != n; d++ {
		for {
			md := m * d
			if n%md != 0 {
				break
			}
			m = md
			c += d - 1
		}
	}
	complexityCache[n] = c
	return c
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

type ratio struct {
	a, b int
}

func (r ratio) float() float64 { return float64(r.a) / float64(r.b) }

var ratios []ratio

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
						ratios = append(ratios, ratio{n, d})
					}
				}
			}
		}
	}
}
