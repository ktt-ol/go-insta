package main

import (
	"flag"
	"math/rand"
	"time"

	"github.com/comail/colog"

	"github.com/ktt-ol/go-insta"
	"github.com/ktt-ol/go-insta/life"
	"github.com/ktt-ol/go-insta/tron"
)

var addrs = []string{
	"192.168.3.6", // # 1 at 0:f:17:10:53:d1
	"192.168.3.8", // # 2 at 0:f:17:10:53:b1
	"192.168.3.9", // # 3 at 0:f:17:10:53:ac
	"192.168.3.7", // # 4 at 0:f:17:10:53:90
	"192.168.3.4", // # 5 at 0:f:17:10:53:c3
	"192.168.3.5", // # 6 at 0:f:17:10:53:b9
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	colog.Register()
	colog.ParseFields(true)

	srv := insta.NewServer()

	var (
		fps     = flag.Int("fps", 25, "fps")
		runLife = flag.Bool("life", false, "life")
		runTron = flag.Bool("tron", false, "tron")
	)

	// c, err := insta.NewInstaClient(addrs)
	// if err != nil {
	//  log.Fatal(err)
	// }

	flag.Parse()

	c := insta.NewTerm()
	c.SetFPS(*fps)
	go c.Run()

	l := life.NewLife(insta.ScreenWidth, insta.ScreenHeight)

	srv.SetLife(l)
	go srv.Run()

	if *runLife {
		s := insta.NewScreen()
		prev := s.Copy()
		sps := 10
		blendSteps := int(float32(*fps) / float32(sps))
		for {
			insta.LifeToScreen(l, s)
			for _, img := range insta.BlendScreens(prev, s, blendSteps) {
				c.SetScreen(img)
				time.Sleep(1000 / time.Duration(*fps) * time.Millisecond)
			}
			l.Step()
			prev, s = s, prev
		}
	}

	if *runTron {
		s := insta.NewScreen()
		tr := tron.NewGame(insta.ScreenWidth, insta.ScreenHeight)
		srv.SetTron(tr)
		for {
			tr.Step()
			tr.Paint(s)
			c.SetScreen(s)
			time.Sleep(1000 / time.Duration(*fps) * time.Millisecond)
		}
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
