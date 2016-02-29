package insta

import (
	"fmt"
	"sync"
	"time"
)

type Term struct {
	screen *Screen
	fps    int
	mu     *sync.Mutex
}

func NewTerm() *Term {
	t := Term{mu: &sync.Mutex{}, fps: 10}
	return &t
}

func (t *Term) SetScreen(s *Screen) {
	t.mu.Lock()
	t.screen = s.Copy()
	t.mu.Unlock()
}

func (t *Term) SetFPS(fps int) {
	t.fps = fps
}

func (t *Term) SetAfterglow(v float64) {
}

const asciiGreyscale = " .'`^\",:;Il!i><~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"

func (t *Term) print() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.screen == nil {
		return
	}

	fmt.Print("\033[2J\033[;H")
	for y := 0; y < t.screen.Bounds().Dy(); y++ {
		for x := 0; x < t.screen.Bounds().Dx(); x++ {
			px := t.screen.At(x, y)
			r, g, b, _ := px.RGBA()
			v := 0.3*float32(r) + 0.6*float32(g) + 0.1*float32(b)
			fmt.Print(string(asciiGreyscale[int(v/0xffff*float32(len(asciiGreyscale)-1))]))
		}
		fmt.Println()
	}

}
func (t *Term) Run() {
	tick := time.Tick(time.Duration(1000.0/float32(t.fps)) * time.Millisecond)
	for _ = range tick {
		t.print()
	}
}
