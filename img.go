package insta

import (
	"image"
	"image/draw"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nfnt/resize"
)

func ShowImage(c Client, fname string) {
	if strings.HasSuffix(fname, ".gif") {
		showGif(c, fname, 0)
		return
	}
	r, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	img, _, err := image.Decode(r)
	if err != nil {
		log.Fatal(err)
	}

	img = resize.Resize(ScreenWidth, ScreenHeight, img, resize.Bilinear)

	scr := NewScreen()
	bnds := image.Rect(0, 0, ScreenWidth, ScreenHeight)
	draw.Draw(scr, bnds, img, image.ZP, draw.Over)

	c.SetScreen(scr)
}

func showGif(c Client, fname string, d time.Duration) {
	r, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	g, err := gif.DecodeAll(r)
	if err != nil {
		log.Fatal(err)
		return
	}
	start := time.Now()
	bnds := image.Rect(0, 0, ScreenWidth, ScreenHeight)
	for {
		scr := NewScreen()
		for i := range g.Image {
			img := resize.Resize(ScreenWidth, ScreenHeight, g.Image[i], resize.Bilinear)
			op := draw.Over
			if g.Disposal[i] == gif.DisposalBackground {
				op = draw.Src
			}
			draw.Draw(scr, bnds, img, image.ZP, op)
			c.SetScreenImmediate(scr)
			// fmt.Println(g.Disposal[i], g.Delay[i], len(g.Delay), g.BackgroundIndex, gif.DisposalPrevious, gif.DisposalBackground, gif.DisposalNone)
			time.Sleep(time.Duration(g.Delay[i]) * time.Millisecond)
		}
		if d == 0 {
			return
		} else if time.Since(start) > d {
			return
		}
	}
}

func RandomGif(c Client, dir string, d time.Duration) {
	gifs, _ := filepath.Glob(filepath.Join(dir, "*.gif"))
	if len(gifs) == 0 {
		return
	}
	i := rand.Intn(len(gifs))
	showGif(c, gifs[i], d)
}

func ScrollImage(c Client, fname string) {
	r, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	img, _, err := image.Decode(r)
	if err != nil {
		log.Fatal(err)
	}

	img = resize.Resize(0, ScreenHeight, img, resize.Bilinear)

	steps := img.Bounds().Dx() + ScreenWidth + 1
	for i := 0; i < steps; i++ {
		scr := NewScreen()
		bnds := image.Rect(ScreenWidth-i, 0, ScreenWidth, ScreenHeight)
		draw.Draw(scr, bnds, img, image.ZP, draw.Over)

		c.SetScreen(scr)
	}
}
