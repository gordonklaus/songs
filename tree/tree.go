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
		sineFreq:     newMelody(256, 16),
		noteDuration: newMelody(4, 16),
	})
}

type song struct {
	EventDelay   audio.EventDelay
	sineFreq     melody
	noteDuration melody
	nodes        []node
	MultiVoice   audio.MultiVoice
}

type node struct {
	childDuration *note
	numChildren   int
}

func (s *song) InitAudio(p audio.Params) {
	audio.Init(&s.EventDelay, p)
	s.EventDelay.Delay(0, s.next)
	audio.Init(&s.MultiVoice, p)
}

func (s *song) next() {
	for {
		n := len(s.nodes)
		if n == 0 {
			break
		}
		s.nodes[n-1].numChildren--
		if s.nodes[n-1].numChildren > 0 {
			break
		}
		s.nodes = s.nodes[:n-1]
	}

	for len(s.nodes) == 0 || len(s.nodes) < 5 {
		s.newNode()
	}

	t := s.nodes[len(s.nodes)-1].childDuration.abs
	s.sineFreq.time += t
	s.noteDuration.time += t
	s.EventDelay.Delay(t, s.next)
}

func (s *song) newNode() {
	var childDuration *note
	numChildren := 1
	if len(s.nodes) > 0 {
		var r ratio
		parent := s.nodes[len(s.nodes)-1]
		childDuration, r = s.noteDuration.nextAfter(parent.childDuration.abs, parent.childDuration, invNatRats)
		numChildren = r.b
	} else {
		childDuration, _ = s.noteDuration.next(0)
		_, r := s.noteDuration.nextAfter(childDuration.abs, childDuration, natRats)
		numChildren = r.a
		childDuration.time.max += float64(numChildren) * childDuration.abs
	}
	s.nodes = append(s.nodes, node{
		childDuration: childDuration,
		numChildren:   numChildren,
	})

	duration := float64(numChildren) * childDuration.abs
	sineFreq, _ := s.sineFreq.next(duration)
	s.MultiVoice.Add(newSineVoice(sineFreq.abs, duration))
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

func (s *song) Sing() float64 {
	s.EventDelay.Step()
	return math.Tanh(s.MultiVoice.Sing() / 8)
}

func (s *song) Done() bool {
	return s.MultiVoice.Done()
}

type sineVoice struct {
	Osc audio.FixedFreqSineOsc
	Env audio.ExpEnv
	amp float64
}

func newSineVoice(freq, duration float64) *sineVoice {
	v := &sineVoice{}
	v.Osc.SetFreq(freq)
	v.Env.Go(1, .1).Go(.7, duration-.2).Go(0, .1)
	v.amp = 4 / math.Log2(freq)
	return v
}

func (v *sineVoice) InitAudio(p audio.Params) {
	v.Osc.InitAudio(p)
	v.Env.InitAudio(p)
}

func (v *sineVoice) Sing() float64 {
	return math.Tanh(2*v.Osc.Sine()) * v.Env.Sing() * v.amp
}

func (v *sineVoice) Done() bool {
	return v.Env.Done()
}
