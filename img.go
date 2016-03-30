package insta

import (
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"time"

	"github.com/nfnt/resize"
)

func ShowImage(c Client, fname string, dur time.Duration) {
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
	syncAt := time.Now()
	for i := 0; i < steps; i++ {
		scr := NewScreen()
		bnds := image.Rect(ScreenWidth-i, 0, ScreenWidth, ScreenHeight)
		draw.Draw(scr, bnds, img, image.ZP, draw.Over)

		// wait till previous frame was synced, in case we are to fast
		if syncAt.After(time.Now()) {
			time.Sleep(syncAt.Sub(time.Now()))
		}

		// step syncAt time for next frame
		syncAt = syncAt.Add(dur)

		// is next frame in the past? forward to next syncAt in the future
		for syncAt.Before(time.Now()) {
			log.Println("dropped frame")
			syncAt = syncAt.Add(dur)
		}
		c.SetScreenAt(scr, syncAt)
	}
}
