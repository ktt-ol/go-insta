package insta

import (
	"image/color"
	"math"
	"math/rand"
	"time"
)

func Rainbow(c Client, d time.Duration) {
	start := time.Now()
	ix := rand.Float64() * 100
	iy := rand.Float64() * 50
	ih := rand.Float64() * 360
	for time.Since(start) < d {
		s := NewScreen()
		for x := 0; x < ScreenWidth; x++ {
			for y := 0; y < ScreenHeight; y++ {
				r, g, b := HsvToRgb(
					float64(y)*1.5+ih,
					0.2+0.8*(math.Sin((ix+float64(x))/ScreenWidth*3)/2+0.5),
					0.1+0.9+(math.Cos((iy+float64(y))/ScreenHeight*3)/2+0.5),
				)
				s.Set(x, y, color.RGBA{uint8(r * 255), uint8(g * 255), uint8(b * 255), 128})
			}
		}
		ix += 0.8
		iy += 0.5
		ih += 2.0
		c.SetScreen(s)
	}
}
