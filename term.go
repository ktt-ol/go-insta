package insta

import (
	"fmt"
	"image"
	"log"
	"time"
)

type Term struct {
	imgs     chan image.Image
	fps      int
	nextSync time.Time
}

func NewTerm() *Term {
	t := Term{imgs: make(chan image.Image, 1), fps: 10}
	return &t
}

func (t *Term) SetScreen(s *Screen) {
	dur := time.Second / time.Duration(t.fps)
	// wait till previous frame was synced, in case we are to fast
	if t.nextSync.After(time.Now()) {
		time.Sleep(t.nextSync.Sub(time.Now()))
	}

	// step t.nextSync time for next frame
	t.nextSync = t.nextSync.Add(dur)

	// is next frame in the past? forward to next t.nextSync in the future
	if t.nextSync.Before(time.Now()) {
		log.Println("dropped frame")
		t.nextSync = time.Now().Add(dur)
	}

	select {
	case t.imgs <- s.Copy():
	default: // skip screen
	}
}

func (t *Term) SetFPS(fps int) {
	t.fps = fps
}

func (t *Term) SetAfterglow(v float64) {
}

const asciiGreyscale = " .'`^\",:;Il!i><~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"

func (t *Term) print(img image.Image) {
	fmt.Print("\033[2J\033[;H")
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			px := img.At(x, y)
			r, g, b, _ := px.RGBA()
			v := 0.3*float32(r) + 0.6*float32(g) + 0.1*float32(b)
			fmt.Print(string(asciiGreyscale[int(v/0xffff*float32(len(asciiGreyscale)-1))]))
		}
		fmt.Println()
	}

}
func (t *Term) Run() {
	for img := range t.imgs {
		t.print(img)
		time.Sleep(time.Duration(1000/float64(t.fps)) * time.Millisecond)
	}
}
