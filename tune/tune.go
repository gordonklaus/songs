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
	Voices voices
	Reverb audio.Reverb `audio:"noscore"`
}

func (b *band) Sing() float64 {
	x := b.Voices.Sing() / 4
	const wet = .3
	return (1-wet)*x + wet*b.Reverb.Filter(x)
}

func (b *band) Done() bool {
	return b.Voices.Done()
}

type voices struct {
	audio.MultiVoice
}

func (s *voices) Play(n struct{ Pitch, Amplitude []*audio.ControlPoint }) {
	s.Add(&voice{
		Pitch: audio.NewControl(n.Pitch),
		Amp:   audio.NewControl(n.Amplitude),
	})
}

func (s *voices) Sing() float64 {
	return s.MultiVoice.Sing()
}

func (s *voices) Done() bool {
	return s.MultiVoice.Done()
}

type voice struct {
	Pitch, Amp *audio.Control
	Saw        audio.SawOsc
	LP         audio.LowPass1
}

func (v *voice) Sing() float64 {
	freq := math.Exp2(v.Pitch.Sing())
	v.Saw.Freq(freq)
	v.LP.Freq(freq)
	return v.LP.Filter(v.Saw.Sing() * math.Exp2(v.Amp.Sing()))
}

func (v *voice) Done() bool {
	return v.Pitch.Done() && v.Amp.Done()
}
