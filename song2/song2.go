package main

import (
	"github.com/gordonklaus/audio"
	"github.com/gordonklaus/audiogui"

	"math"
)

func main() {
	audiogui.Main(score, &band{})
}

type band struct {
	Sines sines
}

func (b *band) Sing() float64 {
	return b.Sines.Sing()
}

func (b *band) Done() bool {
	return b.Sines.Done()
}

type sines struct {
	audio.MultiVoice
}

func (s *sines) Play(n struct{ Pitch, Amplitude []*audio.ControlPoint }) {
	v := &sineVoice{}
	v.Pitch.SetPoints(n.Pitch)
	v.Amp.SetPoints(n.Amplitude)
	d := math.Max(v.Pitch.Duration(), v.Amp.Duration())
	v.Env.AttackHoldRelease(.05, d - .05, 4)
	s.Add(v)
}

type sineVoice struct {
	Pitch, Amp audio.Control
	Env        audio.ExpEnv
	Sine       audio.SineOsc
}

func (v *sineVoice) Sing() float64 {
	f := math.Exp2(v.Pitch.Sing())
	g := math.Exp2(v.Amp.Sing()) * v.Env.Sing()
	return g * math.Tanh(2*v.Sine.Sine(f))
}

func (v *sineVoice) Done() bool {
	return v.Pitch.Done() && v.Amp.Done() && v.Env.Done()
}
