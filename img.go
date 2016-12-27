package insta

import (
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"github.com/nfnt/resize"
)

func ShowImage(c Client, fname string) {
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
