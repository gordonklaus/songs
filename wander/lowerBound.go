package main

import (
	"math"
	"sort"
)

// TODO: b
func getLowerBound(b int, D []int) *lowerBoundIterator {
	offset := D[0]
	D = append([]int{}, D...)
	for i := range D {
		D[i] -= offset
	}

	if len(D) > 2 {
		return newLowerBoundIterator(newOffsetLowerBound(b, offset, D))
	}

	lbc := theLowerBoundCache.get(D)
	if lbc.lb == nil {
		lbc.lb = newLowerBound(b, D)
	}
	return newLowerBoundIterator(&offsetLowerBound{
		offset: offset,
		lb:     lbc.lb,
	})
}

var theLowerBoundCache lowerBoundCache

type lowerBoundCache struct {
	lb   *lowerBound
	next []*lowerBoundCache
}

func (lbc *lowerBoundCache) get(D []int) *lowerBoundCache {
	if len(D) == 0 {
		return lbc
	}
	d := D[0]
	for d >= len(lbc.next) {
		lbc.next = append(lbc.next, nil)
	}
	if lbc.next[d] == nil {
		lbc.next[d] = &lowerBoundCache{}
	}
	return lbc.next[d].get(D[1:])
}

type lowerBoundIterator struct {
	lb               *offsetLowerBound
	i, n0, n1, value int
}

func newLowerBoundIterator(lb *offsetLowerBound) *lowerBoundIterator {
	lbi := &lowerBoundIterator{lb: lb}
	lbi.increment()
	return lbi
}

func (lbi *lowerBoundIterator) increment() {
	lbi.i++
	lb := lbi.lb.lb
	for lbi.i >= lb.pending {
		lb.advance()
	}
	step0 := lb.steps[lbi.i-1]
	step1 := lb.steps[lbi.i]
	lbi.n0 = step0.n - lbi.lb.offset
	lbi.n1 = step1.n - lbi.lb.offset
	lbi.value = step0.value
}

// TODO: merge into lowerBoundIterator
type offsetLowerBound struct {
	offset int
	lb     *lowerBound
}

func newOffsetLowerBound(b, offset int, D []int) *offsetLowerBound {
	return &offsetLowerBound{
		offset: offset,
		lb:     newLowerBound(b, D),
	}
}

type lowerBound struct {
	b int
	D []int

	m       int
	sum     *lowerBoundSum
	steps   []lowerBoundStep
	pending int
}

type lowerBoundStep struct {
	n, value int
}

func newLowerBound(b int, D []int) *lowerBound {
	partials := []*lowerBoundIterator{}
	if PD := partition(D); len(PD) > 1 {
		for _, D := range PD {
			partials = append(partials, getLowerBound(b, D))
		}
	} else if len(D) > 2 {
		for _, d := range D {
			partials = append(partials, getLowerBound(b, []int{d}))
		}
	}

	return &lowerBound{
		b:   b,
		D:   D,
		sum: newLowerBoundSum(partials),
		m:   1,
	}
}

func partition(D []int) [][]int {
	const maxPartitionSize = 2

	if len(D) <= maxPartitionSize {
		return [][]int{D}
	}

	// i := len(D) / 2
	// return append(partition(D[:i]), partition(D[i:])...)

	switch maxPartitionSize {
	case 2:
		// TODO:  No 1-term partitions.
		return append(partition(D[2:]), D[:2])
	case 3:
		if len(D)%3 == 1 {
			return append(partition(D[2:]), D[:2])
		}
		return append(partition(D[3:]), D[:3])
	default:
		panic("not implemented")
	}
}

func (lb *lowerBound) advance() {
	// fmt.Println(lb.D, " --- ", lb.steps[:lb.pending], "-", lb.m)
	// defer func() {
	// 	fmt.Println(lb.D, " --- ", lb.steps[:lb.pending], "-", lb.m)
	// 	fmt.Println()
	// }()

	if len(lb.D) == 1 {
		n := 1
		if lb.pending > 0 {
			n = 1 + 1<<uint(lb.pending-1)
		}
		lb.steps = append(lb.steps, lowerBoundStep{n, lb.pending})
		lb.pending++
		return
	}

	if len(lb.D) == 2 {
		for next := lb.pending + 1; lb.pending < next; {
			advanceTwoTermLowerBounds()
		}
		return
	}

	for ; ; lb.m++ {
		if lb.m >= lb.sum.n {
			if lb.pending < len(lb.steps) && lb.steps[lb.pending].value <= lb.sum.value {
				lb.pending++
				return
			}
			lb.sum.increment()
		}

		if gcd(lb.m, lb.b) != 1 {
			continue
		}
		value := lb.evaluate(lb.m)

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

func (lb *lowerBound) evaluate(n int) int {
	value := 0
	for _, d := range lb.D {
		value += complexity(n + d)
	}
	return value
}

type lowerBoundSum struct {
	lbis     []*lowerBoundIterator
	n, value int
}

func newLowerBoundSum(lbis []*lowerBoundIterator) *lowerBoundSum {
	value := 0
	for _, lbi := range lbis {
		value += lbi.value
	}
	return &lowerBoundSum{lbis, 1, value}
}

func (lbs *lowerBoundSum) increment() {
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

	// for _, lbi := range lbs.lbis {
	// 	fmt.Println(lbi.n0, "--", lbi.n1)
	// }
	// fmt.Println("---")
}

var inverseComplexityCache = [][]int{{1}}

func advanceTwoTermLowerBounds() {
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

				lbc := theLowerBoundCache.get([]int{0, d})
				if lbc.lb == nil {
					lbc.lb = newLowerBound(1, []int{0, d})
				}
				lb := lbc.lb

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
