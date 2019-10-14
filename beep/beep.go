package main

import (
	"fmt"
	"math"
	"sort"
	"sync"

	"github.com/gordonklaus/audio"
	"github.com/gordonklaus/gui"
	"github.com/gordonklaus/portaudio"
)

func main() {
	w := &window{voices: map[int]*sineVoice{}}
	p := audio.Params{SampleRate: 96000}
	audio.Init(w, p)

	portaudio.Initialize()
	defer portaudio.Terminate()
	s, err := portaudio.OpenDefaultStream(0, 1, p.SampleRate, 64, w.processAudio)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := s.Start(); err != nil {
		fmt.Println(err)
		return
	}
	defer s.Stop()

	if err := gui.Run(func() {
		gui.NewWindow(w, "beep", func(win *gui.Window) {
			w.Window = win
			gui.SetKeyFocus(w)
		})
	}); err != nil {
		fmt.Println(err)
	}
}

type window struct {
	*gui.Window

	mu     sync.Mutex
	Multi  audio.MultiVoice
	voices map[int]*sineVoice
}

func (w *window) KeyPress(k gui.KeyEvent) {
	w.mu.Lock()
	defer w.mu.Unlock()

	_, playing := w.voices[k.Key]
	if freq, ok := keyFreq[k.Key]; !playing && ok {
		v := &sineVoice{}
		v.Sine.Freq(freq)
		w.voices[k.Key] = v
		v.Env.AttackHoldRelease(.01, 999, 0)
		w.Multi.Add(v)
	}
}

func (w *window) KeyRelease(k gui.KeyEvent) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if v, ok := w.voices[k.Key]; ok {
		delete(w.voices, k.Key)
		v.Env.ReleaseNow(2)
	}
}

func (w *window) processAudio(out []float32) {
	w.mu.Lock()
	defer w.mu.Unlock()

	for i := range out {
		out[i] = float32(w.Multi.Sing())
	}
}

type sineVoice struct {
	Sine audio.SineOsc
	Env  audio.ExpEnv
}

func (v *sineVoice) Sing() float64 {
	return distort(v.Env.Sing() * v.Sine.Sing())
}

func (v *sineVoice) Done() bool {
	return v.Env.Done()
}

func distort(x float64) float64 {
	y := math.Abs(x) + 1
	return math.Copysign(1-1/(y*y), x)
}

func pitchToFreq(pitch float64) float64 { return 256 * math.Pow(2, pitch/12) }

var keyFreq12 = map[int]float64{
	gui.KeyZ:            pitchToFreq(-12),
	gui.KeyS:            pitchToFreq(-11),
	gui.KeyX:            pitchToFreq(-10),
	gui.KeyD:            pitchToFreq(-9),
	gui.KeyC:            pitchToFreq(-8),
	gui.KeyV:            pitchToFreq(-7),
	gui.KeyG:            pitchToFreq(-6),
	gui.KeyB:            pitchToFreq(-5),
	gui.KeyH:            pitchToFreq(-4),
	gui.KeyN:            pitchToFreq(-3),
	gui.KeyJ:            pitchToFreq(-2),
	gui.KeyM:            pitchToFreq(-1),
	gui.KeyComma:        pitchToFreq(0),
	gui.KeyL:            pitchToFreq(1),
	gui.KeyPeriod:       pitchToFreq(2),
	gui.KeySemicolon:    pitchToFreq(3),
	gui.KeySlash:        pitchToFreq(4),
	gui.KeyQ:            pitchToFreq(0),
	gui.Key2:            pitchToFreq(1),
	gui.KeyW:            pitchToFreq(2),
	gui.Key3:            pitchToFreq(3),
	gui.KeyE:            pitchToFreq(4),
	gui.KeyR:            pitchToFreq(5),
	gui.Key5:            pitchToFreq(6),
	gui.KeyT:            pitchToFreq(7),
	gui.Key6:            pitchToFreq(8),
	gui.KeyY:            pitchToFreq(9),
	gui.Key7:            pitchToFreq(10),
	gui.KeyU:            pitchToFreq(11),
	gui.KeyI:            pitchToFreq(12),
	gui.Key9:            pitchToFreq(13),
	gui.KeyO:            pitchToFreq(14),
	gui.Key0:            pitchToFreq(15),
	gui.KeyP:            pitchToFreq(16),
	gui.KeyLeftBracket:  pitchToFreq(17),
	gui.KeyEqual:        pitchToFreq(18),
	gui.KeyRightBracket: pitchToFreq(19),
	gui.KeyBackslash:    pitchToFreq(21),
}

var keyFreq = map[int]float64{}

func init() {
	xs := []float64{}
	for three := 0.; three <= 2; three++ {
		for five := 0.; five <= 2; five++ {
			for seven := 0.; seven <= 0; seven++ {
				x := math.Pow(3, three) * math.Pow(5, five) * math.Pow(7, seven)
				x /= math.Exp2(math.Trunc(math.Log2(x)))
				xs = append(xs, x)
			}
		}
	}
	sort.Float64s(xs)

	for i, k := range zxcv {
		keyFreq[k] = xs[i%len(xs)] * 128 * math.Exp2(float64(i/len(xs)))
	}
	for i, k := range asdf {
		keyFreq[k] = xs[i%len(xs)] * 256 * math.Exp2(float64(i/len(xs)))
	}
	for i, k := range qwer {
		keyFreq[k] = xs[i%len(xs)] * 512 * math.Exp2(float64(i/len(xs)))
	}
}

var zxcv = [...]int{
	gui.KeyZ,
	gui.KeyX,
	gui.KeyC,
	gui.KeyV,
	gui.KeyB,
	gui.KeyN,
	gui.KeyM,
	gui.KeyComma,
	gui.KeyPeriod,
	gui.KeySlash,
}

var asdf = [...]int{
	gui.KeyA,
	gui.KeyS,
	gui.KeyD,
	gui.KeyF,
	gui.KeyG,
	gui.KeyH,
	gui.KeyJ,
	gui.KeyK,
	gui.KeyL,
	gui.KeySemicolon,
	gui.KeyApostrophe,
	gui.KeyEnter,
}

var qwer = [...]int{
	gui.KeyQ,
	gui.KeyW,
	gui.KeyE,
	gui.KeyR,
	gui.KeyT,
	gui.KeyY,
	gui.KeyU,
	gui.KeyI,
	gui.KeyO,
	gui.KeyP,
	gui.KeyLeftBracket,
	gui.KeyRightBracket,
	gui.KeyBackslash,
}
