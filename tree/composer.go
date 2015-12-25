package main

import "math/rand"

type composer struct {
	frequency melody
	duration  melody
	nodes     []node
	events    chan event
}

type node struct {
	childDuration *note
	numChildren   int
}

type event struct {
	notes []composedNote
	next  float64
}

type composedNote struct {
	frequency, duration float64
}

func newComposer() *composer {
	c := &composer{
		frequency: newMelody(256, 16),
		duration:  newMelody(4, 16),
		events:    make(chan event, 10),
	}
	go c.compose()
	return c
}

func (c *composer) compose() {
	for {
		var e event

		for {
			n := len(c.nodes)
			if n == 0 {
				break
			}
			c.nodes[n-1].numChildren--
			if c.nodes[n-1].numChildren > 0 {
				break
			}
			c.nodes = c.nodes[:n-1]
		}

		for len(c.nodes) == 0 || len(c.nodes) < 9 && rand.Float64() < .95 {
			e.notes = append(e.notes, c.newNode())
		}

		t := c.nodes[len(c.nodes)-1].childDuration.abs
		c.frequency.time += t
		c.duration.time += t
		e.next = t

		c.events <- e
	}
}

func (c *composer) newNode() composedNote {
	var childDuration *note
	numChildren := 1
	if len(c.nodes) > 0 {
		var r ratio
		parent := c.nodes[len(c.nodes)-1]
		childDuration, r = c.duration.nextAfter(parent.childDuration.abs, parent.childDuration, invNatRats)
		numChildren = r.b
	} else {
		childDuration, _ = c.duration.next(0)
		_, r := c.duration.nextAfter(childDuration.abs, childDuration, natRats)
		numChildren = r.a
		childDuration.time.max += float64(numChildren) * childDuration.abs
	}
	c.nodes = append(c.nodes, node{
		childDuration: childDuration,
		numChildren:   numChildren,
	})

	duration := float64(numChildren) * childDuration.abs
	frequency, _ := c.frequency.next(duration)
	return composedNote{
		frequency: frequency.abs,
		duration:  duration,
	}
}

var invNatRats = func() (r []ratio) {
	for i := 1; i <= 6; i++ {
		r = append(r, ratio{1, i})
	}
	return
}()

var natRats = func() (r []ratio) {
	for i := 1; i <= 6; i++ {
		r = append(r, ratio{i, 1})
	}
	return
}()
