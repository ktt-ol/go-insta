package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/ktt-ol/go-insta"
	"github.com/ktt-ol/go-insta/life"
)

var addrs = []string{
	"192.168.3.6", // # 1 at 0:f:17:10:53:d1
	"192.168.3.8", // # 2 at 0:f:17:10:53:b1
	"192.168.3.9", // # 3 at 0:f:17:10:53:ac
	"192.168.3.7", // # 4 at 0:f:17:10:53:90
	"192.168.3.4", // # 5 at 0:f:17:10:53:c3
	"192.168.3.5", // # 6 at 0:f:17:10:53:b9
}

func main() {
	s := insta.NewScreen()

	c, err := insta.NewClient(addrs)
	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().UTC().UnixNano())
	l := life.NewLife(insta.ScreenWidth, insta.ScreenHeight)

	// s.Set(0, 18, color.RGBA{255, 255, 255, 0})
	// fmt.Println(s)
	c.SetScreen(s)

	for {
		l.Step()
		insta.LifeToScreen(l, s)

		if err := c.Send(); err != nil {
			log.Fatal(err)
		}
		time.Sleep(20 * time.Millisecond)
		// return
	}

	// frames := 50
	// blendSteps := 5
	// g := gif.GIF{LoopCount: frames + blendSteps}

	// insta.LifeToScreen(l, s)
	// g.Delay = append(g.Delay, 5 /* *10ms */)
	// g.Image = append(g.Image, insta.ScreenToPalettedImage(s))
	// for i := 0; i < frames; i++ {
	// 	l.Step()
	// 	insta.LifeToScreen(l, s)
	// 	for bs, img := range insta.BlendImages(g.Image[len(g.Image)-1], s, blendSteps) {
	// 		if bs == blendSteps-1 {
	// 			g.Delay = append(g.Delay, 20 /* *10ms */)
	// 		} else {
	// 			g.Delay = append(g.Delay, 5 /* *10ms */)
	// 		}
	// 		g.Image = append(g.Image, img)
	// 	}
	// }

	// f, err := os.Create("/tmp/out.gif")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer f.Close()
	// if err := gif.EncodeAll(f, &g); err != nil {
	// 	log.Fatal(err)
	// }

}
