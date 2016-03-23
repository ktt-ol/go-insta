package insta

import (
	"fmt"
	"image"
	"time"
)

type Term struct {
	imgs chan image.Image
	fps  int
}

func NewTerm() *Term {
	t := Term{imgs: make(chan image.Image, 1), fps: 10}
	return &t
}

func (t *Term) SetScreen(s *Screen) {
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
