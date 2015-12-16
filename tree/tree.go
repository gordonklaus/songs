package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/gordonklaus/audio"
)

var (
	print   = fmt.Print
	printf  = fmt.Printf
	println = fmt.Println
)

func main() {
	rand.Seed(time.Now().UnixNano())
	audio.Play(&song{
		melody:   newMelody(256, 4),
		beatFreq: newMelody(1, 4),
	})
}

type song struct {
	EventDelay audio.EventDelay
	melody     melody
	beatFreq   melody
	beats      []beat
	MultiVoice audio.MultiVoice
	done       bool
}

type beat struct {
	count    int
	duration float64
	sineFreq float64
	note     *note
}

func (s *song) InitAudio(p audio.Params) {
	audio.Init(&s.EventDelay, p)
	s.EventDelay.Delay(0, s.beat)
	audio.Init(&s.MultiVoice, p)
}

func (s *song) beat() {
	for {
		n := len(s.beats)
		if n == 0 {
			break
		}
		s.beats[n-1].count--
		if s.beats[n-1].count > 0 {
			break
		}
		s.beats = s.beats[:n-1]
	}

	if len(s.beats) > 0 {
		s.MultiVoice.Add(newSineVoice(s.beats[len(s.beats)-1]))
	}

	for len(s.beats) == 0 || len(s.beats) < 5 && rand.Float64() < .75 {
		s.newBeat()
	}

	t := 1 / s.beats[len(s.beats)-1].note.f
	s.melody.time += t
	s.beatFreq.time += t
	s.EventDelay.Delay(t, s.beat)
}

func (s *song) newBeat() {
	count := 1
	var note *note
	if len(s.beats) > 0 {
		var r ratio
		note, r = s.beatFreq.nextAfter(s.beats[len(s.beats)-1].note, []ratio{{1, 1}, {2, 1}, {3, 1}, {4, 1}, {5, 1}, {6, 1}})
		count = r.a
	} else {
		note, _ = s.beatFreq.next(allRats)
		count = 1 + rand.Intn(6)
	}
	sineFreq, _ := s.melody.next(allRats)
	b := beat{
		count:    count,
		note:     note,
		sineFreq: sineFreq.f,
	}
	s.beats = append(s.beats, b)
	s.MultiVoice.Add(newSineVoice(b))
}

func (s *song) Sing() float64 {
	s.EventDelay.Step()
	return math.Tanh(s.MultiVoice.Sing() / 8)
}

func (s *song) Done() bool {
	return s.done && s.MultiVoice.Done()
}

type sineVoice struct {
	Osc audio.FixedFreqSineOsc
	Env *audio.AttackReleaseEnv
	amp float64
	n   int
}

func newSineVoice(b beat) *sineVoice {
	v := &sineVoice{}
	v.Osc.SetFreq(b.sineFreq)
	v.Env = audio.NewAttackReleaseEnv(.1, 4)
	v.amp = 4 / math.Log2(b.sineFreq)
	return v
}

func (v *sineVoice) InitAudio(p audio.Params) {
	v.Osc.InitAudio(p)
	v.Env.InitAudio(p)
	v.n = int(p.SampleRate * .1)
}

func (v *sineVoice) Sing() float64 {
	v.n--
	if v.n < 0 {
		v.Env.Release()
	}
	return math.Tanh(2*v.Osc.Sine()) * v.Env.Sing() * v.amp
}

func (v *sineVoice) Done() bool {
	return v.Env.Done()
}

type ratio struct {
	a, b int
}

func (r ratio) float() float64 { return float64(r.a) / float64(r.b) }

type melody struct {
	center        float64
	coherency     float64
	coherencyTime float64

	time    float64
	history []*note
}

type note struct {
	t float64
	f float64
	n int
}

func newMelody(center, coherencyTime float64) melody {
	return melody{
		center:        center,
		coherency:     math.Pow(.01, 1./coherencyTime),
		coherencyTime: coherencyTime,
		history:       []*note{{0, center, 1}},
	}
}

func (m *melody) next(rats []ratio) (*note, ratio) {
	return m.nextAfter(m.history[len(m.history)-1], rats)
}

func (m *melody) nextAfter(prev *note, rats []ratio) (*note, ratio) {
	cSum, ampSum := m.historyComplexity()

	sum := 0.0
	sums := make([]float64, len(rats))
	for i, r := range rats {
		p := math.Log2(prev.f * r.float() / m.center)
		sum += math.Exp2(-p*p/2) * math.Exp2(-m.complexity(prev, cSum, ampSum, r))
		sums[i] = sum
	}
	i := sort.SearchFloat64s(sums, sum * rand.Float64())
	next := m.appendHistory(prev, rats[i])

	for i, n := range m.history {
		if m.time-n.t < m.coherencyTime {
			m.history = m.history[i:]
			d := m.history[0].n
			for _, n := range m.history[1:] {
				d = gcd(d, n.n)
			}
			for i := range m.history {
				m.history[i].n /= d
			}
			break
		}
	}

	return next, rats[i]
}

var allRats []ratio

func init() {
	pow := func(a, x int) int {
		y := 1
		for x > 0 {
			y *= a
			x--
		}
		return y
	}
	mul := func(n, d, a, x int) (int, int) {
		if x > 0 {
			return n * pow(a, x), d
		}
		return n, d * pow(a, -x)
	}
	for _, two := range []int{-3, -2, -1, 0, 1, 2, 3} {
		for _, three := range []int{-2, -1, 0, 1, 2} {
			for _, five := range []int{-1, 0, 1} {
				for _, seven := range []int{-1, 0, 1} {
					n, d := 1, 1
					n, d = mul(n, d, 2, two)
					n, d = mul(n, d, 3, three)
					n, d = mul(n, d, 5, five)
					n, d = mul(n, d, 7, seven)
					if complexity(n, d) < 12 {
						allRats = append(allRats, ratio{n, d})
					}
				}
			}
		}
	}
}

func (m *melody) historyComplexity() (cSum, ampSum float64) {
	for i, n1 := range m.history {
		a1 := math.Pow(m.coherency, m.time-n1.t)
		for _, n2 := range m.history[:i] {
			a2 := math.Pow(m.coherency, m.time-n2.t)
			cSum += a1 * a2 * float64(complexity(n1.n, n2.n))
		}
		ampSum += a1
	}
	return
}

func (m *melody) complexity(prev *note, cSum, ampSum float64, r ratio) float64 {
	const a1 = 1
	n1n := r.a * prev.n
	for _, n2 := range m.history {
		a2 := math.Pow(m.coherency, m.time-n2.t)
		cSum += a1 * a2 * float64(complexity(n1n, n2.n*r.b))
	}
	return cSum / (ampSum + a1)
}

func complexity(a, b int) int {
	c := 0
	for d := 2; a != b; {
		d1 := a%d == 0
		d2 := b%d == 0
		if d1 != d2 {
			c += d - 1
		}
		if d1 {
			a /= d
		}
		if d2 {
			b /= d
		}
		if !(d1 || d2) {
			d++
		}
	}
	return c
}

func (m *melody) appendHistory(prev *note, r ratio) *note {
	prevN := prev.n
	for i := range m.history {
		m.history[i].n *= r.b
	}
	n := &note{m.time, prev.f * r.float(), r.a * prevN}
	m.history = append(m.history, n)
	return n
}

func gcd(a, b int) int {
	for a > 0 {
		if a > b {
			a, b = b, a
		}
		b -= a
	}
	return b
}
