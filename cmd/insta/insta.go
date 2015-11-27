package main

import (
	"flag"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/comail/colog"
	"github.com/tarm/serial"

	"github.com/ktt-ol/go-insta"
	"github.com/ktt-ol/go-insta/life"
	"github.com/ktt-ol/go-insta/tron"
)

var addrs = []string{
	"192.168.3.6", // # 1 at 0:f:17:10:53:d1
	"192.168.3.9", // # 2 at 0:f:17:10:53:ac
	"192.168.3.8", // # 3 at 0:f:17:10:53:b1
	"192.168.3.5", // # 4 at 0:f:17:10:53:b9
	"192.168.3.7", // # 5 at 0:f:17:10:53:90
	"192.168.3.4", // # 6 at 0:f:17:10:53:c3
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	colog.Register()
	colog.ParseFields(true)

	var (
		fps      = flag.Int("fps", 25, "fps")
		runLife  = flag.Bool("life", false, "life")
		runPaint = flag.Bool("paint", false, "paint")
		runTron  = flag.Bool("tron", false, "tron")
		port     = flag.String("port", "", "serial port")
	)

	c, err := insta.NewInstaClient(addrs)
	if err != nil {
		log.Fatal(err)
	}
	// c := insta.NewTerm()

	flag.Parse()

	var ser *serial.Port
	if *port != "" {
		cfg := &serial.Config{}
		cfg.Baud = 57600
		cfg.Name = *port
		var err error
		ser, err = serial.OpenPort(cfg)
		if err != nil {
			log.Fatal(err)
		}
	}

	mp := insta.NewMultiPad(ser)
	mp.Pads()

	c.SetFPS(*fps)
	go c.Run()

	l := life.NewLife(insta.ScreenWidth, insta.ScreenHeight)

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
		for {
			tr.Step(mp.Pads())
			tr.Paint(s)
			c.SetScreen(s)
			time.Sleep(1000 / time.Duration(*fps) * time.Millisecond)
		}
	}

	if *runPaint {
		s := insta.NewScreen()
		for {
			p := mp.Pads()[0]
			x, y := p.StickLeft()
			// fmt.Println(x, y, int((x+1)/2*insta.ScreenWidth), int((y+1)/2*insta.ScreenHeight))
			s.Set(int((x+1)/2*insta.ScreenWidth), int((y+1)/2*insta.ScreenHeight), color.RGBA{120, 4, 200, 0})
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
