package main

import (
	"fmt"
)

func getLowerBound(b, offset int, D []int) *lowerBoundIterator {
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
	lbi := &lowerBoundIterator{lb: lb, i: -1}
	lbi.increment()
	return lbi
}

func (lbi *lowerBoundIterator) increment() {
	lbi.i++
	lb := lbi.lb.lb
	if lbi.i >= len(lb.steps) {
		lb.advance()
	}
	lbi.n0 = lbi.n1
	lbi.n1 = lb.steps[lbi.i].n - lbi.lb.offset
	lbi.value = lb.steps[lbi.i].value
}

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
	b        int
	D        []int
	partials []*lowerBoundIterator

	l, m    int
	pending []lowerBoundStep
	steps   []lowerBoundStep
}

type lowerBoundStep struct {
	n, value int
}

func newLowerBound(b int, D []int) *lowerBound {
	partials := []*lowerBoundIterator{}
	if len(D) > 1 {
		// for i := 1; i < len(D); i += 2 {
		// 	d0 := D[i-1]
		// 	d1 := D[i]
		// 	partials = append(partials, getLowerBound(b, d0, []int{0, d1 - d0}))
		// }
		// if len(D)%2 == 1 {
		// 	partials = append(partials, getLowerBound(b, D[len(D)-1], []int{0}))
		// }
		for _, d := range D {
			partials = append(partials, getLowerBound(b, d, []int{0}))
		}
	}
	return &lowerBound{
		b:        b,
		D:        D,
		partials: partials,
		m:        1,
	}
}

func (lb *lowerBound) advance() {
	for ; lb.m < lb.l || len(lb.pending) < 2; lb.m++ {
		if gcd(lb.m, lb.b) != 1 {
			continue
		}
		value := lb.evaluate(lb.m)

		i := len(lb.pending)
		for ; i > 0; i-- {
			if value > lb.pending[i-1].value {
				break
			}
		}
		if i < len(lb.pending) {
			lb.pending = lb.pending[:i+1]
			lb.pending[i].value = value
		} else {
			lb.pending = append(lb.pending, lowerBoundStep{lb.m, value})
		}
		if i == 0 {
			lb.updateLimit()
		}
	}

	lb.steps = append(lb.steps, lowerBoundStep{lb.pending[1].n, lb.pending[0].value})
	lb.pending = lb.pending[1:]
	lb.updateLimit()
}

func (lb *lowerBound) evaluate(n int) int {
	value := 0
	for _, d := range lb.D {
		value += complexity(n + d)
	}
	return value
}

func (lb *lowerBound) updateLimit() {
	value := lb.pending[0].value

	if len(lb.D) == 1 {
		lb.l = 1 << uint(value)
		return
	}

	for {
		weaklb := 0
		maxN0 := 0
		plbMinN1 := lb.partials[0]
		for _, plb := range lb.partials {
			weaklb += plb.value
			if plb.n0 > maxN0 {
				maxN0 = plb.n0
			}
			if plb.n1 < plbMinN1.n1 {
				plbMinN1 = plb
			}
		}
		if weaklb > value {
			lb.l = maxN0
			break
		}
		plbMinN1.increment()
	}

	fmt.Println("-", lb.m, lb.l, lb.pending)
}
