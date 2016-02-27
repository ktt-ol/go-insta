package insta

import (
	"image/color"
	"math"
	"math/rand"
	"time"
)

type star struct {
	x     float64
	y     float64
	dx    float64
	dy    float64
	alive bool
	hue   int
}

func Spaceflight(c Client, duration time.Duration) {
	stars := make([]star, 400)

	till := time.Now().Add(duration)
	for {
		scr := NewScreen()
		added := 0
		alive := 0
		for i, s := range stars {
			if added < 2 && !s.alive && time.Now().Before(till) {
				s = star{
					x:     ScreenWidth / 2,
					y:     ScreenHeight / 2,
					dx:    (rand.Float64()*2 + 0.2) - 1.1,
					dy:    (rand.Float64()*2 + 0.2) - 1.1,
					alive: true,
					hue:   rand.Intn(360),
				}
				added += 1
			}
			if !s.alive {
				continue
			}
			alive += 1
			s.x += s.dx
			s.y += s.dy
			s.dx *= 1.05
			s.dy *= 1.05
			if s.x < 0 || s.y < 0 || s.x >= ScreenWidth || s.y >= ScreenHeight {
				s.alive = false
			} else {
				dist := math.Hypot((s.x-ScreenWidth/2)/ScreenWidth/2, (s.y-ScreenHeight/2)/ScreenHeight/2)
				saturation := dist * 3
				if saturation > 1.0 {
					saturation = 1.0
				}
				brightness := 0.1 + dist*3
				if brightness > 1.0 {
					brightness = 1.0
				}

				r, g, b := hsvToRgb(float64(s.hue), saturation, brightness)
				scr.Set(int(s.x), int(s.y), color.RGBA{uint8(r * 255), uint8(g * 255), uint8(b * 255), 128})
			}
			stars[i] = s
		}
		c.SetScreen(scr)
		time.Sleep(40 * time.Millisecond)
		if alive == 0 {
			break
		}
	}
}
