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
	nextDuration  []ratioComplexity
	nextFrequency []ratioComplexity
}

type note struct {
	t, f ratio
}

type ratioComplexity struct {
	r ratio
	c int
}

func NewMelody() *Melody {
	rhythmComplexity := .8    // 0..1
	frequencyComplexity := .5 // 0..1
	avgDuration := 0.25
	avgFrequency := 256.0
	coherencyTime := 8.0
	m := &Melody{
		rhythmBias:    math.Log2(rhythmComplexity),
		frequencyBias: math.Log2(frequencyComplexity),
		avgDuration:   avgDuration,
		avgFrequency:  avgFrequency,
		coherencyTime: coherencyTime,
		lastDuration:  avgDuration,
		lastFrequency: avgFrequency,
	}
	for i := 0; i < int(coherencyTime/avgDuration); i++ {
		m.appendHistory(ratio{1, 1}, ratio{1, 1})
	}
	return m
}

func (m *Melody) Next() (float64, float64) {
	m.appendHistory(m.newDuration(), m.newFrequency())
	return m.lastDuration, m.lastFrequency
}

func (m *Melody) newDuration() ratio {
	return selectRatio(m.nextDuration, func(rc ratioComplexity) float64 {
		n := len(m.history)
		n = n * (n + 1) / 2
		n = n * (n - 1) / len(m.history)
		c := float64(rc.c) / float64(n)
		return math.Exp(-m.lastDuration*rc.r.float()/m.avgDuration) * math.Exp2(m.rhythmBias*c)
	})
}

func (m *Melody) newFrequency() ratio {
	return selectRatio(m.nextFrequency, func(rc ratioComplexity) float64 {
		n := len(m.history)
		n = n * (n + 1) / len(m.history)
		c := float64(rc.c) / float64(n)
		dp := math.Log2(m.lastFrequency * rc.r.float() / m.avgFrequency)
		return math.Exp2(-dp*dp/2) * math.Exp2(m.frequencyBias*c)
	})
}

func selectRatio(candidates []ratioComplexity, weight func(ratioComplexity) float64) ratio {
	sum := 0.0
	sums := make([]float64, len(candidates))
	for i, rc := range candidates {
		sum += weight(rc)
		sums[i] = sum
	}
	return candidates[sort.SearchFloat64s(sums, sum*rand.Float64())].r
}

func (m *Melody) appendHistory(rd, rf ratio) {
	for i, dc := range m.nextDuration {
		if !dc.r.lessThan(rd) && dc.r != rd {
			m.nextDuration = m.nextDuration[i:]
			break
		}
	}

	histlen := len(m.history) - 2
	if histlen < 0 {
		histlen = 0
	}
	for _, n := range m.history[:histlen] {
		if (rd.float()-n.t.float())*m.lastDuration <= m.coherencyTime {
			break
		}

		for i := range m.nextDuration {
			dc := &m.nextDuration[i]
			dc.c -= m.firstDurationComplexity(dc.r)
		}
		for i := range m.nextFrequency {
			fc := &m.nextFrequency[i]
			fc.c -= m.firstFrequencyComplexity(fc.r)
		}

		m.history = m.history[1:]
	}

	m.history = append(m.history, note{
		t: rd,
		f: rf,
	})
	for i := range m.history {
		n := &m.history[i]
		n.t = n.t.sub(rd).div(rd).normalized()
		n.f = n.f.div(rf).normalized()
	}
	m.lastDuration *= rd.float()
	m.lastFrequency *= rf.float()
	// fmt.Println(rd, "---", m.history)

	for i := range m.nextDuration {
		dc := &m.nextDuration[i]
		dc.r = dc.r.sub(rd).div(rd).normalized()
	}
	for i := range m.nextFrequency {
		fc := &m.nextFrequency[i]
		fc.r = fc.r.div(rf).normalized()
	}

	for i := range m.nextDuration {
		dc := &m.nextDuration[i]
		dc.c += m.nextDurationComplexity(dc.r)
	}
	for i := range m.nextFrequency {
		fc := &m.nextFrequency[i]
		fc.c += m.nextFrequencyComplexity(fc.r)
	}

	// for _, dc := range m.nextDuration {
	// 	if c := m.durationComplexity(dc.r); c != dc.c {
	// 		print(c-dc.c, ", ")
	// 	} else {
	// 		// print("|")
	// 	}
	// }
	m.nextDuration = addNext(m.nextDuration, m.durationComplexity)
	m.nextFrequency = addNext(m.nextFrequency, m.frequencyComplexity)

	// fmt.Println(len(m.nextDuration), len(m.nextFrequency))
}

func addNext(rcs []ratioComplexity, complexity func(ratio) int) []ratioComplexity {
	ret := []ratioComplexity{}
	ir := 0
	for _, rc := range rcs {
		for ; ir < len(ratios) && ratios[ir].lessThan(rc.r); ir++ {
			r := ratios[ir]
			ret = append(ret, ratioComplexity{r, complexity(r)})
		}
		if ir < len(ratios) && ratios[ir] == rc.r {
			ir++
		}
		ret = append(ret, rc)
	}
	for ; ir < len(ratios); ir++ {
		r := ratios[ir]
		ret = append(ret, ratioComplexity{r, complexity(r)})
	}
	return ret
}

func (m *Melody) firstDurationComplexity(next ratio) int {
	c := 0
	d := next.sub(m.history[0].t)
	for i, n1 := range m.history {
		for _, n0 := range m.history[:i] {
			r := n1.t.sub(n0.t).div(d)
			c += complexity(r.a, r.b)
		}
	}
	n0 := m.history[0]
	for _, n1 := range m.history[1:] { // Start at 1 to avoid d1 == 0.
		d1 := n1.t.sub(n0.t)
		for _, n2 := range m.history[1:] { // Start at 1 because we already counted the first one above.
			d2 := next.sub(n2.t)
			r := d2.div(d1)
			c += complexity(r.a, r.b)
		}
	}
	return c
}

func (m *Melody) firstFrequencyComplexity(next ratio) int {
	r := next.div(m.history[0].f)
	return complexity(r.a, r.b)
}

func (m *Melody) durationComplexity(next ratio) int {
	c := 0
	for i, n1 := range m.history {
		for _, n0 := range m.history[:i] {
			d1 := n1.t.sub(n0.t)
			for _, n2 := range m.history {
				d2 := next.sub(n2.t)
				r := d2.div(d1)
				c += complexity(r.a, r.b)
			}
		}
	}
	return c
}

func (m *Melody) frequencyComplexity(next ratio) int {
	c := 0
	for _, n := range m.history {
		r := next.div(n.f)
		c += complexity(r.a, r.b)
	}
	return c
}

func (m *Melody) nextDurationComplexity(next ratio) int {
	c := 0
	for i, n1 := range m.history {
		for _, n0 := range m.history[:i] {
			r := n1.t.sub(n0.t).div(next)
			c += complexity(r.a, r.b)
		}
	}
	for _, n1 := range m.history[:len(m.history)-1] { // Stop at len-1 to avoid d1 == 0.
		d1 := ratio{0, 1}.sub(n1.t)
		for _, n2 := range m.history[:len(m.history)-1] { // Stop at len-1 because we already counted the last one above.
			d2 := next.sub(n2.t)
			r := d2.div(d1)
			c += complexity(r.a, r.b)
		}
	}
	return c
}

func (m *Melody) nextFrequencyComplexity(next ratio) int {
	return complexity(next.a, next.b)
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
	if a > b {
		a, b = b, a
	}
	for a > 0 {
		a, b = b%a, a
	}
	return b
}

type ratio struct {
	a, b int
}

func (r ratio) normalized() ratio {
	neg := false
	if r.a < 0 {
		r.a = -r.a
		neg = !neg
	}
	if r.b < 0 {
		r.b = -r.b
		neg = !neg
	}
	d := gcd(r.a, r.b)
	r.a /= d
	r.b /= d
	if neg {
		r.a = -r.a
	}
	return r
}

func (r ratio) sub(s ratio) ratio {
	d := r.b * s.b
	return ratio{r.a*s.b - s.a*r.b, d}
}

func (r ratio) div(s ratio) ratio {
	return ratio{r.a * s.b, r.b * s.a}
}

func (r ratio) lessThan(s ratio) bool {
	return r.a*s.b < r.b*s.a
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
	sort.Sort(ratiosAscending(ratios))
}

type ratiosAscending []ratio

func (r ratiosAscending) Len() int           { return len(r) }
func (r ratiosAscending) Less(i, j int) bool { return r[i].lessThan(r[j]) }
func (r ratiosAscending) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
