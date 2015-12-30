package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/gordonklaus/audio"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	audio.Play(&song{composer: newComposer()})
}

type song struct {
	composer   *composer
	EventDelay audio.EventDelay
	MultiVoice audio.MultiVoice
}

func (s *song) InitAudio(p audio.Params) {
	audio.Init(&s.EventDelay, p)
	s.EventDelay.Delay(0, s.next)
	audio.Init(&s.MultiVoice, p)
}

func (s *song) next() {
	select {
	case e := <-s.composer.events:
		for _, n := range e.notes {
			s.MultiVoice.Add(newSineVoice(n))
		}
		s.EventDelay.Delay(e.next, s.next)
	default:
		fmt.Print(".")
		s.EventDelay.Delay(.01, s.next)
	}
}

func (s *song) Sing() float64 {
	s.EventDelay.Step()
	return audio.Saturate(s.MultiVoice.Sing() / 4)
}

func (s *song) Done() bool {
	return false
}

type sineVoice struct {
	Osc audio.SawOsc
	LP  audio.LowPass1
	Env audio.ExpEnv
	amp float64
}

func newSineVoice(n composedNote) *sineVoice {
	v := &sineVoice{}
	v.Osc.Freq(n.frequency)
	v.LP.Freq(n.frequency)
	v.Env.Go(1, .1).Go(.7, n.duration-.2).Go(0, .1)
	v.amp = 4 / math.Log2(n.frequency)
	return v
}

func (v *sineVoice) Sing() float64 {
	return v.LP.Filter(v.Osc.Sing()) * v.Env.Sing() * v.amp
}

func (v *sineVoice) Done() bool {
	return v.Env.Done()
}
