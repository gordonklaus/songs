package main

import (
	"math"
	"sort"
)

// TODO: b
func getLowerBoundA(b int, D []int) *lowerBoundIterator {
	D = append([]int{}, D...)
	sort.Ints(D)
	offset := D[0]
	for i := range D {
		D[i] -= offset
	}

	if len(D) == 2 {
		d := D[1]
		for d >= len(lowerBoundACache) {
			lowerBoundACache = append(lowerBoundACache, &lowerBoundA{})
		}
		return newLowerBoundIterator(offset, lowerBoundACache[d])
	}

	lbis := make([]*lowerBoundIterator, len(D))
	for i := range D {
		lbis[i] = getLowerBoundA(b, []int{D[i], D[(i+1)%len(D)]})
	}
	eval := func(n int) int {
		// Fetch all complexities in range, for speed.
		// complexity(n + D[len(D)-1])
		// nums := numbers[n : n+D[len(D)-1]+1]
		value := 0
		for _, d := range D {
			value += complexity(n + d)
		}
		return value
	}
	return newLowerBoundIterator(offset, newLowerBoundSum(b, lbis, true, eval))
}

func getLowerBoundB(D []int) *lowerBoundIterator {
	D = append([]int{}, D...)
	sort.Ints(D)
	if D[0] != 0 {
		panic("expected zero")
	}
	D = D[1:]

	lbis := make([]*lowerBoundIterator, len(D))
	for i, d := range D {
		for d >= len(lowerBoundBCache) {
			lowerBoundBCache = append(lowerBoundBCache, &lowerBoundB{})
		}
		lbis[i] = newLowerBoundIterator(0, lowerBoundBCache[d])
	}
	eval := func(n int) int {
		value := 0
		for _, d := range D {
			value += complexity(n*d + 1)
		}
		return value
	}
	return newLowerBoundIterator(0, newLowerBoundSum(1, lbis, false, eval))
}

type lowerBoundIterator struct {
	offset           int
	lb               lowerBound
	i, n0, n1, value int
}

func newLowerBoundIterator(offset int, lb lowerBound) *lowerBoundIterator {
	return &lowerBoundIterator{offset: offset, lb: lb}
}

func (lbi *lowerBoundIterator) increment() {
	lbi.i++
	for lbi.i >= len(lbi.lb.getSteps()) {
		lbi.lb.advance()
	}
	steps := lbi.lb.getSteps()
	step0 := steps[lbi.i-1]
	step1 := steps[lbi.i]
	lbi.n0 = step0.n - lbi.offset
	lbi.n1 = step1.n - lbi.offset
	lbi.value = step0.value
}

type lowerBound interface {
	getSteps() []lowerBoundStep
	advance()
}

type lowerBoundSum struct {
	b int

	m        int
	n, value int
	lbis     []*lowerBoundIterator
	double   bool
	eval     func(int) int

	steps   []lowerBoundStep
	pending int
}

type lowerBoundStep struct {
	n, value int
}

func newLowerBoundSum(b int, lbis []*lowerBoundIterator, double bool, eval func(int) int) *lowerBoundSum {
	return &lowerBoundSum{
		b:      b,
		m:      1,
		n:      1,
		lbis:   lbis,
		double: double,
		eval:   eval,
	}
}

func (lbs *lowerBoundSum) getSteps() []lowerBoundStep { return lbs.steps[:lbs.pending] }

func (lb *lowerBoundSum) advance() {
	for ; ; lb.m++ {
		// fmt.Println("lowerBoundSum.advance", lb.m)
		if lb.m >= lb.n {
			mul := 1
			if lb.double {
				mul = 2
			}
			if lb.pending < len(lb.steps) && lb.steps[lb.pending].value*mul <= lb.value {
				lb.pending++
				return
			}
			lb.incrementSum()
		}

		// if gcd(lb.m, lb.b) != 1 {
		// 	continue
		// }
		value := lb.eval(lb.m)

		i := len(lb.steps)
		for ; i >= lb.pending && i > 0; i-- {
			if value > lb.steps[i-1].value {
				break
			}
		}
		if i < len(lb.steps) {
			lb.steps = lb.steps[:i+1]
			lb.steps[i].value = value
		} else {
			lb.steps = append(lb.steps, lowerBoundStep{lb.m, value})
		}
	}
}

func (lbs *lowerBoundSum) incrementSum() {
	lbiMinN1 := lbs.lbis[0]
	for _, lbi := range lbs.lbis {
		if lbi.n1 < lbiMinN1.n1 {
			lbiMinN1 = lbi
		}
	}
	lbs.n = lbiMinN1.n1
	lbs.value -= lbiMinN1.value
	lbiMinN1.increment()
	lbs.value += lbiMinN1.value
}

type lowerBoundA struct {
	steps   []lowerBoundStep
	pending int
}

func (lb *lowerBoundA) getSteps() []lowerBoundStep { return lb.steps[:lb.pending] }

func (lb *lowerBoundA) advance() {
	for next := lb.pending + 1; lb.pending < next; {
		advanceLowerBoundAs()
	}
}

var lowerBoundAComplexity int
var lowerBoundACache []*lowerBoundA

func advanceLowerBoundAs() {
	lowerBoundAComplexity++
	c := lowerBoundAComplexity

	for c >= len(inverseComplexityCache) {
		advanceInverseComplexities()
	}

	for c1, c2 := 0, c; c1 <= c2; c1, c2 = c1+1, c2-1 {
		for _, n1 := range inverseComplexityCache[c1] {
			for _, n2 := range inverseComplexityCache[c2] {
				n := n1 + 1
				d := n2 - n1
				if n2 < n1 {
					n = n2 + 1
					d = n1 - n2
				}

				if d > 128 {
					if n1 < n2 {
						break
					}
					continue
				}

				for d >= len(lowerBoundACache) {
					lowerBoundACache = append(lowerBoundACache, &lowerBoundA{})
				}
				lb := lowerBoundACache[d]

				if len(lb.steps) == 0 {
					lb.steps = []lowerBoundStep{{1, 0}}
				}

				numSteps := len(lb.steps)
				lastStep := &lb.steps[numSteps-1]
				var lastStep2 *lowerBoundStep
				if numSteps > 1 {
					lastStep2 = &lb.steps[numSteps-2]
				}
				if lastStep2 == nil || lastStep2.value < c && lastStep.n < n {
					lastStep.value = c
					lb.steps = append(lb.steps, lowerBoundStep{n, 0})
					lb.pending++
				} else if lastStep2.value == c && lastStep.n < n {
					lastStep.n = n
				}
			}
		}
	}
}

type lowerBoundB struct {
	steps   []lowerBoundStep
	pending int
}

func (lb *lowerBoundB) getSteps() []lowerBoundStep { return lb.steps[:lb.pending] }

func (lb *lowerBoundB) advance() {
	for next := lb.pending + 1; lb.pending < next; {
		advanceLowerBoundBs()
	}
}

var lowerBoundBComplexity int
var lowerBoundBCache []*lowerBoundB

func advanceLowerBoundBs() {
	lowerBoundBComplexity++
	c := lowerBoundBComplexity

	// fmt.Println("advanceLowerBoundBs:", c)

	for c >= len(inverseComplexityCache) {
		advanceInverseComplexities()
	}

	for _, n := range inverseComplexityCache[c] {
		for d, lb := range lowerBoundBCache {
			if d == 0 || (n-1)%d != 0 {
				continue
			}

			if len(lb.steps) == 0 {
				lb.steps = []lowerBoundStep{{1, 0}}
			}

			n := (n-1)/d + 1

			numSteps := len(lb.steps)
			lastStep := &lb.steps[numSteps-1]
			var lastStep2 *lowerBoundStep
			if numSteps > 1 {
				lastStep2 = &lb.steps[numSteps-2]
			}
			if lastStep2 == nil || lastStep2.value < c && lastStep.n < n {
				lastStep.value = c
				lb.steps = append(lb.steps, lowerBoundStep{n, 0})
				lb.pending++
			} else if lastStep2.value == c && lastStep.n < n {
				lastStep.n = n
			}
		}
	}
}

var inverseComplexityCache = [][]int{{1}}

func advanceInverseComplexities() {
	c := len(inverseComplexityCache)

	ics := []int{}
	for i := 0; ; i++ {
		p := prime(i)
		if p-1 > c {
			break
		}
		if c%(p-1) != 0 {
			continue
		}
		for _, n := range inverseComplexityCache[c-p+1] {
			if n > math.MaxInt64/p {
				break
			}
			ics = append(ics, n*p)
		}
	}
	inverseComplexityCache = append(inverseComplexityCache, uniqueSort(ics))
}

func uniqueSort(s []int) []int {
	sort.Ints(s)
	i := 0
	for j := 0; j < len(s); i++ {
		s[i] = s[j]
		for j < len(s) && s[i] == s[j] {
			j++
		}
	}
	return s[:i]
}
