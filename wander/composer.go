package main

type composer struct {
	melody *Melody
	notes  chan composedNote
}

type composedNote struct {
	duration, frequency float64
}

func newComposer() *composer {
	c := &composer{
		melody: NewMelody(),
		notes:  make(chan composedNote, 40),
	}
	go c.compose()
	return c
}

func (c *composer) compose() {
	for {
		d, f := c.melody.Next()
		c.notes <- composedNote{d, f}
	}
}
