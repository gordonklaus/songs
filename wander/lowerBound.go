package main

func getLowerBound(b, offset int, D []int) *lowerBoundIterator {
	// if len(D) > 2 {
	// 	return newLowerBoundIterator(newOffsetLowerBound(b, offset, D))
	// }
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
	if lbi.i >= lb.pending || lbi.i >= len(lb.steps) {
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
	b   int
	D   []int
	sum *lowerBoundSum

	m       int
	steps   []lowerBoundStep
	pending int
}

type lowerBoundStep struct {
	n, value int
}

func newLowerBound(b int, D []int) *lowerBound {
	partials := []*lowerBoundIterator{}
	// if len(D) == 2 {
	// 	for _, d := range D {
	// 		partials = append(partials, getLowerBound(b, d, []int{0}))
	// 	}
	// }
	if len(D) >= 2 {
		i := len(D) / 2
		D2 := make([]int, len(D)-i)
		for j := range D2 {
			D2[j] = D[i+j] - D[i]
		}
		partials = append(
			partials,
			getLowerBound(b, 0, D[:i]),
			getLowerBound(b, D[i], D2),
		)
		// for i := 1; i < len(D); i += 2 {
		// 	d0 := D[i-1]
		// 	d1 := D[i]
		// 	partials = append(partials, getLowerBound(b, d0, []int{0, d1 - d0}))
		// }
		// if len(D)%2 == 1 {
		// 	partials = append(partials, getLowerBound(b, D[len(D)-1], []int{0}))
		// }
		// for _, d := range D {
		// 	partials = append(partials, getLowerBound(b, d, []int{0}))
		// }
	}
	lb := &lowerBound{
		b:   b,
		D:   D,
		sum: newLowerBoundSum(partials),
		m:   1,
	}
	lb.advance()
	return lb
}

func (lb *lowerBound) advance() {
	if len(lb.D) == 1 {
		n := 1
		if lb.pending > 0 {
			n = 1 + 1<<uint(lb.pending-1)
		}
		lb.steps = append(lb.steps, lowerBoundStep{n, lb.pending})
		lb.pending++
		return
	}

	for ; ; lb.m++ {
		// fmt.Println(lb.D, "advance", lb.m)

		if lb.m >= lb.sum.n {
			if lb.pending < len(lb.steps) && lb.steps[lb.pending].value <= lb.sum.value {
				break
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
		// fmt.Println("steps:", lb.steps, lb.pending)
	}

	lb.pending++
	// fmt.Println(lb.D, "advanced", lb.steps, lb.pending)
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
}
