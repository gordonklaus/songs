package main

type composer struct {
	frequency melody
	duration  melody
	notes     chan composedNote
}

type composedNote struct {
	frequency, duration float64
}

func newComposer() *composer {
	c := &composer{
		frequency: newMelody(256, 32),
		duration:  newMelody(.5, 32),
		notes:     make(chan composedNote, 10),
	}
	go c.compose()
	return c
}

func (c *composer) compose() {
	for {
		f := c.frequency.next()
		d := c.duration.next()
		c.frequency.time += d
		c.duration.time += d
		c.notes <- composedNote{f, d}
	}
}
