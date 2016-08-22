package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/gordonklaus/audio"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	audio.Play(&song{
		composer: newComposer(),
		Reverb:   audio.NewReverb(),
	})
}

type song struct {
	EventDelay audio.EventDelay
	composer   *composer
	MultiVoice audio.MultiVoice
	Reverb     *audio.Reverb
}

func (s *song) InitAudio(p audio.Params) {
	audio.Init(&s.EventDelay, p)
	s.EventDelay.Delay(0, s.beat)
	audio.Init(&s.MultiVoice, p)
	audio.Init(&s.Reverb, p)
}

func (s *song) beat() {
	n := <-s.composer.notes
	s.MultiVoice.Add(newSineVoice(n.frequency))
	s.EventDelay.Delay(n.duration, s.beat)
}

func (s *song) Sing() float64 {
	s.EventDelay.Step()
	x := audio.Saturate(s.MultiVoice.Sing())
	return audio.Saturate((2*x + s.Reverb.Filter(x)) / 3)
}

func (s *song) Done() bool {
	return false
}

type sineVoice struct {
	Osc audio.SineOsc
	Env audio.ExpEnv
	amp float64
}

func newSineVoice(freq float64) *sineVoice {
	v := &sineVoice{}
	v.Osc.Freq(freq)
	v.Env.AttackHoldRelease(.1, 0, 1)
	v.amp = 4 / math.Log2(freq)
	return v
}

func (v *sineVoice) InitAudio(p audio.Params) {
	v.Osc.InitAudio(p)
	v.Env.InitAudio(p)
}

func (v *sineVoice) Sing() float64 {
	return audio.Saturate(2*v.Osc.Sing()) * v.Env.Sing() * v.amp
}

func (v *sineVoice) Done() bool {
	return v.Env.Done()
}
