package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
)

var _ = fmt.Println

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

	times []timeMultiplicity
}

type note struct {
	t, f ratio
}

type ratioComplexity struct {
	r ratio
	c int
}

type timeMultiplicity struct {
	t ratio
	m int
}

func NewMelody() *Melody {
	rhythmComplexity := .8    // 0..1
	frequencyComplexity := .5 // 0..1
	avgDuration := 1.0
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
	rd := m.newDuration()
	m.appendHistory(rd, m.newFrequency())
	if rd.a == 0 {
		return 0, m.lastFrequency
	}
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
	m.times = nil
	var prev *timeMultiplicity
	for i, n := range m.history {
		if i != 1 && prev != nil && n.t == prev.t {
			prev.m++
		} else {
			m.times = append(m.times, timeMultiplicity{n.t, 1})
			prev = &m.times[len(m.times)-1]
		}
	}

	m.trimPastDurations(rd)

	for i, n := range m.times {
		if i >= len(m.times)-2 || (rd.float()-n.t.float())*m.lastDuration <= m.coherencyTime {
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
		m.times[0].m--
		if m.times[0].m == 0 {
			m.times = m.times[1:]
		}
	}

	m.history = append(m.history, note{
		t: rd,
		f: rf,
	})
	for i := range m.history {
		n := &m.history[i]
		if rd.a > 0 {
			n.t = n.t.sub(rd).div(rd).normalized()
		}
		n.f = n.f.div(rf).normalized()
	}
	if rd.a > 0 {
		m.lastDuration *= rd.float()
	}
	m.lastFrequency *= rf.float()
	// fmt.Println(rd, "---", m.history)

	if rd.a > 0 {
		for i := range m.nextDuration {
			dc := &m.nextDuration[i]
			dc.r = dc.r.sub(rd).div(rd).normalized()
		}
	}
	for i := range m.nextFrequency {
		fc := &m.nextFrequency[i]
		fc.r = fc.r.div(rf).normalized()
	}

	m.times = nil
	prev = nil
	for i, n := range m.history {
		if i != len(m.history)-1 && prev != nil && n.t == prev.t {
			prev.m++
		} else {
			m.times = append(m.times, timeMultiplicity{n.t, 1})
			prev = &m.times[len(m.times)-1]
		}
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
	fmt.Println(rd)
	fmt.Println(m.history)
	{
		oldRatios := ratios
		ratios = m.genNextDurations()
		m.nextDuration = addNext(m.nextDuration, m.durationComplexity)
		ratios = oldRatios
	}
	m.nextFrequency = addNext(m.nextFrequency, m.frequencyComplexity)
	fmt.Println(m.nextDuration)
	fmt.Println()

	// mc := m.newMinComplexity()
	// c, cmia, cmi := mc.minComplexity(mc.lcm, 1)
	// fmt.Println(m.durationComplexity(ratio{1, 1}), c, cmia, cmi)
	// fmt.Println()

	// n := len(m.history)
	// n = n * (n + 1) / 2
	// n = n * (n - 1) / len(m.history)
	// mc := m.newMinComplexity()
	// for _, rc := range m.nextDuration {
	// 	if rc.r.a != 0 {
	// 		r := rc.r
	// 		r.a *= mc.lcm
	// 		r = r.normalized()
	// 		if (float64(rc.c)-mc.minComplexity(r.a, r.b))/float64(n) < 0 {
	// 			panic("no!")
	// 		}
	// 		// fmt.Println(rc, (float64(rc.c)-mc.minComplexity(r.a, r.b))/float64(n))
	// 	}
	// }
	// fmt.Println()
}

func (m *Melody) genNextDurations() []ratio {
	if len(m.history) <= 2 {
		return nil
	}

	minComplexity := 0.0
	minComplexityRatio := ratio{0, 1}
	for i, rc := range m.nextDuration {
		if i == 0 || float64(rc.c) < minComplexity {
			minComplexity = float64(rc.c)
			minComplexityRatio = rc.r
		}
	}
	n := len(m.history)
	n = n * (n + 1) / 2
	n = n * (n - 1) / len(m.history)
	minComplexity /= float64(n)
	pmax := math.Exp(-m.lastDuration*minComplexityRatio.float()/m.avgDuration) * math.Exp2(m.rhythmBias*minComplexity)
	// fmt.Println("most likely:", minComplexityRatio, pmax)

	nextDurations := []ratio{}

	mc := m.newMinComplexity()
	amax := 0
	amaxsum := 0
	for b := 1; ; b++ {
		// fmt.Print("  b:", b)
		for a := 1; ; a++ {
			if gcd(a, b) != 1 {
				continue
			}
			// fmt.Println("  a:", a)

			r := ratio{a, b * mc.lcm}.normalized()
			durationMultiplier := math.Exp(-m.lastDuration * r.float() / m.avgDuration)

			const plimit = .01

			if !(a == 1 && b == 1) {
				cnda := mc.estimateNonDecreasingWithA(a, b)
				pnda := durationMultiplier * math.Exp2(m.rhythmBias*cnda/float64(n))
				if pnda/pmax < plimit {
					if a == 1 {
						cnd := mc.estimateNonDecreasing(a, b)
						pnd := durationMultiplier * math.Exp2(m.rhythmBias*cnd/float64(n))
						if pnd/pmax < plimit {
							fmt.Println("max b:", b, "   max a:", amax, "   max a avg:", float64(amaxsum)/float64(b))
							sort.Sort(ratiosAscending(nextDurations))
							return nextDurations
						}
						// fmt.Println("-a:", a, " ", cnd, " ", pnd/pmax)
					}
					if a > amax {
						amax = a
					}
					amaxsum += a
					break
				}
			}

			c := mc.estimate(a, b)
			p := durationMultiplier * math.Exp2(m.rhythmBias*c/float64(n))
			// fmt.Println("p:", p)
			if p/pmax >= plimit {
				nextDurations = append(nextDurations, r)
			}
		}
		// fmt.Println()
	}
}

type minComplexity struct {
	history   []int
	D         float64
	divCounts []divCount
	coeff     []float64
	GD        float64
	lcm       int
}

type divCount struct {
	mindiv, mindivComplexity int
	G                        float64
}

func (m *Melody) newMinComplexity() minComplexity {
	history := make([]int, len(m.history))
	lcm_ := 1
	for _, n := range m.history {
		lcm_ = lcm(lcm_, n.t.b)
	}
	for i, n := range m.history {
		history[i] = n.t.mul(ratio{lcm_, 1}).normalized().a
	}

	D := 0.0
	for i, t1 := range history {
		for _, t0 := range history[:i] {
			if t1 > t0 {
				D += float64(complexity(t1 - t0))
			}
		}
	}

	divCounts := []divCount{}
	coeff := make([]float64, len(history))
	for i := len(m.history) - 1; i >= 0 && m.history[i].t.a == 0; i-- {
		coeff[i] = 1
	}
	GD := 0.0
	for i := 0; ; i++ {
		p := prime(i)

		tdiv := 0
		for i, t := range history {
			if (-t)%p != 0 {
				tdiv++
				if coeff[i] == 0 {
					coeff[i] = float64(p-1) / math.Log2(float64(p))
				}
			}
		}

		// Delay this check so that all the coeffs can be set above.
		if p > -history[0] {
			break
		}

		mindiv := 1
		mindivComplexity := 0
		G := 0.0
		for d := p; ; d *= p {
			r := make([]int, d)
			for _, t := range history {
				r[(-t)%d]++
			}
			min := len(history)
			max := 0
			for _, r := range r {
				if r < min {
					min = r
				}
				if r > max {
					max = r
				}
			}

			count := 0
			for i, t1 := range history {
				for _, t0 := range history[:i] {
					if t1 > t0 && (t1-t0)%d == 0 {
						count++
					}
				}
			}

			if min == 0 && count == 0 {
				break
			}

			for x := 0; x < min; x++ {
				mindiv *= p
			}
			mindivComplexity += min * (p - 1)
			G += float64((max + 1) * count * (p - 1))
			GD += float64((tdiv + 1) * count * (p - 1))
		}
		divCounts = append(divCounts, divCount{mindiv, mindivComplexity, G})
	}
	for i := range coeff {
		if coeff[i] == 0 {
			panic("ZERO")
		}
	}

	return minComplexity{
		history:   history,
		D:         D,
		divCounts: divCounts,
		coeff:     coeff,
		GD:        GD,
		lcm:       lcm_,
	}
}

func (mc minComplexity) estimate(a, b int) float64 {
	G := 0.0
	for i, divCount := range mc.divCounts {
		p := primes[i]
		if b%p == 0 { // || a%p == 0 ?
			continue
		}
		G += divCount.G
	}

	T := 0
	for _, t := range mc.history {
		T += complexity(a - t*b)
	}

	B := float64(complexity(b))

	N := float64(len(mc.history))
	return (N+2)*(N-1)/2*float64(T) + N*mc.D + N*N*(N-1)/2*B - 2*G
}

func (mc minComplexity) estimateNonDecreasingWithA(a, b int) float64 {
	mindiv := 1
	mindivComplexity := 0
	G := 0.0
	for i, divCount := range mc.divCounts {
		p := primes[i]
		if b%p == 0 {
			continue
		}
		mindiv *= divCount.mindiv
		mindivComplexity += divCount.mindivComplexity
		G += divCount.G
	}

	A := 1.0
	for _, t := range mc.history {
		A *= float64(a - t*b)
	}
	if math.IsInf(A, 0) {
		panic("too much")
	}
	T := math.Log2(A/float64(mindiv)) + float64(mindivComplexity) // consider maxdiv (can it be nondecreasing?)
	B := float64(complexity(b))

	N := float64(len(mc.history))
	return (N+2)*(N-1)/2*T + N*mc.D + N*N*(N-1)/2*B - 2*G
}

func (mc minComplexity) estimateNonDecreasing(a, b int) float64 {
	T := 0.0
	for i, t := range mc.history {
		T += math.Log2(float64(a-t*b)) * mc.coeff[i]
	}
	if math.IsInf(T, 0) {
		panic("too much")
	}
	B := math.Log2(float64(b))

	N := float64(len(mc.history))
	G := float64(mc.GD)
	return (N+2)*(N-1)/2*T + N*mc.D + N*N*(N-1)/2*B - 2*G
}

// func (mc minComplexity) minComplexity(a, b int) (float64, float64, float64) {
// 	ac := 0.0
// 	dc := 0.0
// 	mindiv := 1
// 	maxdiv := []int{}
// 	for i := 0; ; i++ {
// 		p := prime(i)
// 		if p > len(mc.history) && p > a-mc.history[0]*b {
// 			break
// 		}
// 		maxdiv_ := 0
// 		// d <= len(mc.history) is for mindiv   TODO: should be len(mc.times)?
// 		// d <= a-mc.history[0]*b is for maxdiv
// 		for di, d := 0, p; d <= len(mc.history) || d <= a-mc.history[0]*b; d *= p {
// 			if b%p == 0 {
// 				break
// 			}
// 			r := make([]int, d)
// 			for _, t := range mc.history {
// 				r[(-t)%d]++
// 			}
// 			min := d // TODO: len(mc.history) (should make no difference, but for clarity)
// 			max := 0
// 			for _, r := range r {
// 				if r < min {
// 					min = r
// 				}
// 				if r > max {
// 					max = r
// 				}
// 			}
// 			for x := 0; x < min; x++ {
// 				mindiv *= p
// 			}
// 			ac += float64(min * (p - 1))
// 			if i < len(mc.divCounts) && di < len(mc.divCounts[i]) {
// 				dc += float64(max * mc.divCounts[i][di] * (p - 1))
// 			}
// 			maxdiv_ += max - min
// 			di++
// 		}
// 		maxdiv = append(maxdiv, maxdiv_)
// 	}

// 	A := 1.0
// 	// AA := 1
// 	for _, t := range mc.history {
// 		A *= float64(a - t*b)
// 		// AA *= a - t*b
// 	}
// 	if math.IsInf(A, 0) {
// 		panic("too much")
// 	}
// 	// if AA%mindiv != 0 {
// 	// 	fmt.Println(a, b, mc.history)
// 	// 	println(AA, mindiv)
// 	// 	panic("UHOH")
// 	// }
// 	// fmt.Print(math.Log2(A), " ")
// 	A /= float64(mindiv)
// 	// fmt.Println(math.Log2(A) + float64(ac))

// 	for i, maxdiv_ := range maxdiv {
// 		p := float64(primes[i])
// 		d := 1.0
// 		for j := 0; j < maxdiv_; j++ {
// 			d *= p
// 		}
// 		if A < float64(d) {
// 			A = 1
// 			ac += (p - 1) * math.Log(A) / math.Log(p)
// 			break
// 		}
// 		A /= d
// 		ac += float64(maxdiv_ * (primes[i] - 1))
// 	}
// 	for i := len(maxdiv); A > 1; i++ {
// 		p := float64(prime(i))
// 		if A < p {
// 			A = 1
// 			ac += (p - 1) * math.Log(A) / math.Log(p)
// 			break
// 		}
// 		A /= p
// 		ac += p - 1
// 	}

// 	acnm := 0
// 	for _, t := range mc.history {
// 		acnm += complexity(a - t*b)
// 	}

// 	bnm := float64(complexity(b))
// 	bmi := math.Log2(float64(b))

// 	N := float64(len(mc.history))
// 	c := (N+2)*(N-1)/2*float64(acnm) - 2*dc - mc.D + N*N*(N-1)/2*bnm
// 	cmia := (N+2)*(N-1)/2*ac - 2*dc - mc.D + N*N*(N-1)/2*bnm
// 	cmi := (N+2)*(N-1)/2*ac - 2*dc - mc.D + N*N*(N-1)/2*bmi
// 	return c, cmia, cmi
// }

var primes = []int{2, 3}

func prime(i int) int {
	for i >= len(primes) {
		for x := primes[len(primes)-1] + 2; ; x += 2 {
			isPrime := true
			for _, p := range primes {
				if x%p == 0 {
					isPrime = false
					break
				}
			}
			if isPrime {
				primes = append(primes, x)
				break
			}
		}
	}
	return primes[i]
}

func (m *Melody) trimPastDurations(rd ratio) {
	if len(m.times) == 0 {
		return
	}
	lastSimultaneous := 1
	if rd.a == 0 {
		lastSimultaneous = 1 + m.times[len(m.times)-1].m
	}
	const maxSimultaneous = 5

	for i, dc := range m.nextDuration {
		if rd.lessThan(dc.r) || (dc.r == rd && lastSimultaneous < maxSimultaneous) {
			m.nextDuration = m.nextDuration[i:]
			break
		}
	}
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
	d := next.sub(m.times[0].t)
	for i, n1 := range m.times {
		for _, n0 := range m.times[:i] {
			r := n1.t.sub(n0.t).div(d)
			c += n1.m * n0.m * r.complexity()
		}
	}
	n0 := m.times[0]
	for _, n1 := range m.times {
		d1 := n1.t.sub(n0.t)
		for _, n2 := range m.times[1:] { // Start at 1 because we already counted the first one above.
			d2 := next.sub(n2.t)
			r := d2.div(d1)
			c += n1.m * n2.m * r.complexity()
		}
	}
	for _, n1 := range m.times[1:] {
		r := next.sub(n0.t).div(next.sub(n1.t))
		c += n1.m * r.complexity()
	}
	return c
}

func (m *Melody) firstFrequencyComplexity(next ratio) int {
	r := next.div(m.history[0].f)
	return r.complexity()
}

func (m *Melody) durationComplexity(next ratio) int {
	c := 0
	for i, n1 := range m.times {
		for _, n0 := range m.times[:i] {
			d1 := n1.t.sub(n0.t)
			m1 := n1.m * n0.m
			for _, n2 := range m.times {
				d2 := next.sub(n2.t)
				r := d2.div(d1)
				c += m1 * n2.m * r.complexity()
			}

			r := next.sub(n0.t).div(next.sub(n1.t))
			c += m1 * r.complexity()
		}
	}
	return c
}

func (m *Melody) frequencyComplexity(next ratio) int {
	c := 0
	for _, n := range m.history {
		r := next.div(n.f)
		c += r.complexity()
	}
	return c
}

func (m *Melody) nextDurationComplexity(next ratio) int {
	c := 0
	for i, n1 := range m.times {
		for _, n0 := range m.times[:i] {
			r := n1.t.sub(n0.t).div(next)
			c += n1.m * n0.m * r.complexity()
		}
	}
	for _, n1 := range m.times {
		d1 := ratio{0, 1}.sub(n1.t)
		for _, n2 := range m.times[:len(m.times)-1] { // Stop at len-1 because we already counted the last one above.
			d2 := next.sub(n2.t)
			r := d2.div(d1)
			c += n1.m * n2.m * r.complexity()
		}
	}
	for _, n0 := range m.times[:len(m.times)-1] {
		r := next.sub(n0.t).div(next)
		c += n0.m * r.complexity()
	}
	return c
}

func (m *Melody) nextFrequencyComplexity(next ratio) int {
	return next.complexity()
}

var complexityCache = map[int]int{}

func complexity(n int) int {
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
	// if a < 0 {
	// 	a = -a
	// }
	// if b < 0 {
	// 	b = -b
	// }
	if a > b {
		a, b = b, a
	}
	for a > 0 {
		a, b = b%a, a
	}
	return b
}

func lcm(a, b int) int {
	return a * b / gcd(a, b)
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

func (r ratio) mul(s ratio) ratio {
	return ratio{r.a * s.a, r.b * s.b}
}

func (r ratio) div(s ratio) ratio {
	return ratio{r.a * s.b, r.b * s.a}
}

func (r ratio) lessThan(s ratio) bool {
	return r.a*s.b < r.b*s.a
}

func (r ratio) complexity() int {
	if r.a == r.b || r.a == 0 || r.b == 0 {
		return 0
	}
	d := gcd(r.a, r.b)
	d *= d
	return complexity(r.a * r.b / d)
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
					if (ratio{n, d}).complexity() < 12 {
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
