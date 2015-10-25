package main

import (
	"image/color"
	"image/gif"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/ktt-ol/go-insta"
	"github.com/ktt-ol/go-insta/life"
)

func main() {
	s := insta.NewScreen()
	s.Set(0, 0, color.RGBA{255, 32, 128, 0})
	s.Set(17, 17, color.RGBA{255, 32, 128, 0})
	// fmt.Println(s)

	// c, err := insta.NewClient([]string{"192.168.1.1:9410"})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// c.SetScreen(s)

	// if err := c.Send(); err != nil {
	// 	log.Fatal(err)
	// }

	rand.Seed(time.Now().UTC().UnixNano())
	l := life.NewLife(insta.ScreenWidth, insta.ScreenHeight)

	frames := 50
	blendSteps := 5
	g := gif.GIF{LoopCount: frames + blendSteps}

	insta.LifeToScreen(l, s)
	g.Delay = append(g.Delay, 5 /* *10ms */)
	g.Image = append(g.Image, insta.ScreenToPalettedImage(s))
	for i := 0; i < frames; i++ {
		l.Step()
		insta.LifeToScreen(l, s)
		for bs, img := range insta.BlendImages(g.Image[len(g.Image)-1], s, blendSteps) {
			if bs == blendSteps-1 {
				g.Delay = append(g.Delay, 20 /* *10ms */)
			} else {
				g.Delay = append(g.Delay, 5 /* *10ms */)
			}
			g.Image = append(g.Image, img)
		}
	}

	f, err := os.Create("/tmp/out.gif")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := gif.EncodeAll(f, &g); err != nil {
		log.Fatal(err)
	}

}
