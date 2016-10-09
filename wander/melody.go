package main

import (
	"math"
	"math/rand"
)

// Melody generates a melody.
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

// NewMelody creates a Melody.
func NewMelody() *Melody {
	rhythmComplexity := .85   // 0..1
	frequencyComplexity := .5 // 0..1
	avgRate := 1.0
	avgFrequency := 256.0
	return &Melody{
		rhythmBias:    math.Log2(rhythmComplexity),
		frequencyBias: math.Log2(frequencyComplexity),
		avgDuration:   1 / avgRate,
		avgFrequency:  avgFrequency,
		coherencyTime: 8,
		lastDuration:  1 / avgRate,
		lastFrequency: avgFrequency,
		history:       []note{{0, 1}, {1, 1}, {2, 1}},
	}
}

// Next generates a new note and returns its duration and frequency.
func (m *Melody) Next() (float64, float64) {
	rd := m.newDuration()
	rf := m.newFrequency()
	m.appendHistory(rd, rf)

	d, f := m.lastDuration, m.lastFrequency
	if rd.a == 0 {
		d = 0
	}
	return d, f
}

type harmony struct {
	t interval
	n int
}

type interval struct {
	t0, t1 int
}

func (i interval) overlaps(j interval) bool {
	return true //i.t0 < j.t1 && j.t0 < i.t1
}

func (m *Melody) newDuration() ratio {
	// harmonies := make([]harmony, len(m.history)-1)
	// for i := range harmonies {
	// 	harmonies[i] = harmony{
	// 		t: interval{0, 1},
	// 		n: m.history[i+1].t - m.history[i].t,
	// 	}
	// }
	harmonies := make([]harmony, 0, len(m.history)*(len(m.history)-1)/2)
	for i, n1 := range m.history {
		for _, n0 := range m.history[:i] {
			harmonies = append(harmonies, harmony{
				t: interval{n0.t, n1.t},
				n: n1.t - n0.t,
			})
		}
	}
	// times := make([]int, len(m.history))
	// for i, n := range m.history {
	// 	times[i] = n.t
	// }
	// cSum := m.complexitySum3(times)

	cSum := 1 //m.complexitySum(harmonies)

	rats := ratios
	// rats = append(rats, ratio{0, 1})
	sum := 0.0
	sums := make([]float64, len(rats))
	for i, r := range rats {
		f := r.float()
		sum += math.Exp(-m.lastDuration*f/m.avgDuration) * math.Exp2(m.rhythmBias*m.durationComplexity(harmonies, cSum, r))
		// sum += math.Exp(-m.lastDuration*f/m.avgDuration) * math.Exp2(m.rhythmBias*m.durationComplexity3(times, cSum, r))
		sums[i] = sum
	}
	i := 0
	x := sum * rand.Float64()
	for i = range sums {
		if x < sums[i] {
			break
		}
	}

	if rats[i].a != 0 {
		m.lastDuration *= rats[i].float()
	}
	return rats[i]
}

func (m *Melody) newFrequency() ratio {
	harmonies := make([]harmony, len(m.history))
	for i, n := range m.history {
		harmonies[i] = harmony{
			t: interval{0, 1},
			n: n.f,
		}
	}

	cSum := 1 //m.complexitySum(harmonies)

	rats := ratios
	sum := 0.0
	sums := make([]float64, len(rats))
	for i, r := range rats {
		f := r.float()
		dp := math.Log2(m.lastFrequency * f / m.avgFrequency)
		sum += math.Exp2(-dp*dp/2) * math.Exp2(m.frequencyBias*m.frequencyComplexity(harmonies, cSum, r))
		sums[i] = sum
	}
	i := 0
	x := sum * rand.Float64()
	for i = range sums {
		if x < sums[i] {
			break
		}
	}

	m.lastFrequency *= rats[i].float()
	return rats[i]
}

func (m *Melody) complexitySum(harmonies []harmony) int {
	c := 0
	for i, h1 := range harmonies {
		for _, h2 := range harmonies[:i] {
			if h1.t.overlaps(h2.t) {
				c += complexity(h1.n, h2.n)
			}
		}
	}
	return c
}

func (m *Melody) complexitySum3(times []int) int {
	c := 0
	for i1, t1 := range times {
		for i2, t2 := range times[i1+1:] {
			for _, t3 := range times[i1+i2+2:] {
				c += complexity3(t2-t1, t3-t2, t3-t1)
			}
		}
	}
	return c
}

func (m *Melody) durationComplexity(harmonies []harmony, cSum int, r ratio) float64 {
	// last := 1
	// for i := len(m.history) - 1; i >= 0; i-- {
	// 	if m.history[i].n != 0 {
	// 		last = m.history[i].n
	// 		break
	// 	}
	// }
	harmonies2 := make([]harmony, len(m.history))
	last := m.history[len(m.history)-1]
	last2 := m.history[len(m.history)-2]
	d := r.a * (last.t - last2.t)
	t1 := last.t*r.b + d
	for i, n := range m.history {
		t0 := n.t * r.b
		harmonies2[i] = harmony{
			t: interval{t0, t1},
			n: t1 - t0,
		}
	}
	cSum = m.complexity(harmonies, harmonies2, cSum, r)
	n := len(harmonies) + len(harmonies2)
	n = n * (n - 1) / len(m.history)
	return float64(cSum) / float64(n)
}

func (m *Melody) durationComplexity3(times []int, cSum int, r ratio) float64 {
	// last := 1
	// for i := len(m.history) - 1; i >= 0; i-- {
	// 	if m.history[i].n != 0 {
	// 		last = m.history[i].n
	// 		break
	// 	}
	// }
	last := times[len(times)-1]
	last2 := times[len(times)-2]
	d := r.a * (last - last2)
	t3 := last*r.b + d
	for i1, t1 := range times {
		for _, t2 := range times[i1+1:] {
			cSum += complexity3(t3-t2*r.b, (t2-t1)*r.b, t3-t1*r.b)
		}
	}
	n := len(times) + 1
	n = n * (n - 1) / 3
	return float64(cSum) / float64(n)
}

func (m *Melody) frequencyComplexity(harmonies []harmony, cSum int, r ratio) float64 {
	n := r.a * harmonies[len(harmonies)-1].n
	harmonies2 := []harmony{{t: interval{0, 1}, n: n}}
	cSum = m.complexity(harmonies, harmonies2, cSum, r)
	n = len(harmonies) + len(harmonies2)
	n = n * (n - 1) / len(m.history)
	return float64(cSum) / float64(n)
}

func (m *Melody) complexity(harmonies, harmonies2 []harmony, cSum int, r ratio) int {
	for i, h1 := range harmonies2 {
		for _, h2 := range harmonies2[:i] {
			if h1.t.overlaps(h2.t) {
				cSum += complexity(h1.n, h2.n)
			}
		}
	}
	for _, h1 := range harmonies {
		t1 := h1.t
		if true { // should only be done for rhythm, but ok because intervals are hacked to (0, 1) for frequencies
			t1.t0 *= r.b
			t1.t1 *= r.b
		}
		for _, h2 := range harmonies2 {
			if t1.overlaps(h2.t) {
				cSum += complexity(h1.n*r.b, h2.n)
			}
		}
	}
	return cSum
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

func complexity3(a, b, z int) int {
	if a == 0 || b == 0 || z == 0 {
		return 0
	}
	c := 0
	for d := 2; a != b && a != z && b != z; {
		d1 := a%d == 0
		d2 := b%d == 0
		d3 := z%d == 0
		if (d1 || d2 || d3) && !(d1 && d2 && d3) {
			c += d - 1
		}
		if d1 {
			a /= d
		}
		if d2 {
			b /= d
		}
		if d3 {
			z /= d
		}
		if !(d1 || d2 || d3) {
			d++
		}
	}
	return c
}

func (m *Melody) appendHistory(rd, rf ratio) {
	// last := 1
	// for i := len(m.history) - 1; i >= 0; i-- {
	// 	if m.history[i].n != 0 {
	// 		last = m.history[i].n
	// 		break
	// 	}
	// }
	last := m.history[len(m.history)-1]  // Copy the last note before it is modified.
	last2 := m.history[len(m.history)-2] // Copy the last note before it is modified.
	for i := range m.history {
		m.history[i].t *= rd.b
		m.history[i].f *= rf.b
	}
	d := rd.a * (last.t - last2.t)
	t1 := last.t*rd.b + d
	m.history = append(m.history, note{
		t: t1,
		f: rf.a * last.f,
	})

	for i, n := range m.history[:len(m.history)-2] {
		r := ratio{t1 - n.t, d}
		if m.lastDuration*r.float() < m.coherencyTime {
			m.history = m.history[i:]
			break
		}
	}

	div := m.history[0]
	t0 := div.t
	div.t = 0
	for i := range m.history {
		n := &m.history[i]
		n.t -= t0
		div.t = gcd(div.t, n.t)
		div.f = gcd(div.f, n.f)
	}
	for i := range m.history {
		m.history[i].t /= div.t
		m.history[i].f /= div.f
	}
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
