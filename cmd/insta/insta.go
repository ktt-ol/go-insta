package main

import (
	"flag"
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
	"172.16.42.101", // 1 at 0:f:17:10:53:d1
	"172.16.42.102", // 2 at 0:f:17:10:53:ac
	"172.16.42.103", // 3 at 0:f:17:10:53:b1
	"172.16.42.104", // 4 at 0:f:17:10:53:b9
	"172.16.42.105", // 5 at 0:f:17:10:53:90
	"172.16.42.106", // 6 at 0:f:17:10:53:c3
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	colog.Register()
	colog.ParseFields(true)

	var (
		fps     = flag.Int("fps", 25, "fps")
		runLife = flag.Bool("life", false, "life")
		runTron = flag.Bool("tron", false, "tron")
		port    = flag.String("port", "", "serial port")
		term    = flag.Bool("term", false, "use terminal")
	)

	flag.Parse()

	var (
		c   insta.Client
		err error
	)
	if *term {
		log.Println("using terminal")
		c = insta.NewTerm()
	} else {
		log.Println("connecting to", addrs)
		c, err = insta.NewInstaClient(addrs)
		if err != nil {
			log.Fatal(err)
		}
	}

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

	if *runLife {
		l := life.NewLife(insta.ScreenWidth, insta.ScreenHeight)
		go func() {
			t := time.Tick(2 * time.Second)
			for _ = range t {
				l.AddRandomSpaceship()
			}
		}()

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
}
