package main

import (
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/tarm/serial"

	"github.com/ktt-ol/go-insta"
	"github.com/ktt-ol/go-insta/life"
	"github.com/ktt-ol/go-insta/snake"
	"github.com/ktt-ol/go-insta/tron"
)

var addrs = []string{
	"192.168.3.101", // 1 at 0:f:17:10:53:d1
	"192.168.3.102", // 2 at 0:f:17:10:53:ac
	"192.168.3.103", // 3 at 0:f:17:10:53:b1
	"192.168.3.104", // 4 at 0:f:17:10:53:b9
	"192.168.3.105", // 5 at 0:f:17:10:53:90
	"192.168.3.106", // 6 at 0:f:17:10:53:c3

	// "172.16.42.101", // 1 at 0:f:17:10:53:d1
	// "172.16.42.102", // 2 at 0:f:17:10:53:ac
	// "172.16.42.103", // 3 at 0:f:17:10:53:b1
	// "172.16.42.104", // 4 at 0:f:17:10:53:b9
	// "172.16.42.105", // 5 at 0:f:17:10:53:90
	// "172.16.42.106", // 6 at 0:f:17:10:53:c3
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	var (
		fps            = flag.Int("fps", 25, "fps")
		runLife        = flag.Duration("life", 0, "run life for duration")
		runSnake       = flag.Duration("snake", 0, "run snake for duration")
		runTron        = flag.Duration("tron", 0, "run tron for duration")
		runSpaceflight = flag.Duration("spaceflight", 0, "run spaceflight for duration")
		runLogo        = flag.Bool("logo", false, "show mainframe logo")
		port           = flag.String("port", "", "serial port")
		term           = flag.Bool("term", false, "use terminal")
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
			log.Fatal("unable to connect", err)
		}
	}

	var ser *serial.Port

	var pads func() []insta.Pad

	if *port != "" {
		cfg := &serial.Config{}
		cfg.Baud = 57600
		cfg.Name = *port
		var err error
		ser, err = serial.OpenPort(cfg)
		if err != nil {
			log.Fatal(err)
		}
		mp := insta.NewMultiPad(ser)
		pads = mp.Pads
	} else if *term {
		kp := insta.NewKeyboardPad()
		pads = kp.Pads
	} else {
		pads = func() []insta.Pad {
			return nil
		}
	}

	c.SetFPS(*fps)
	go c.Run()

	for {
		if *runLogo {
			c.SetAfterglow(0)
			insta.ShowImage(c, "img/mainframe-mod.png", 25*time.Millisecond) // time.Second/time.Duration(*fps))
			c.SetAfterglow(0.4)
		}

		if runLife.Seconds() > 0 {
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

			till := time.Now().Add(*runLife)
			for time.Now().Before(till) {
				insta.LifeToScreen(l, s)
				for _, img := range insta.BlendScreens(prev, s, blendSteps) {
					c.SetScreen(img)
					time.Sleep(1000 / time.Duration(*fps) * time.Millisecond)
				}
				l.Step()
				prev, s = s, prev
			}
		}

		if runSpaceflight.Seconds() > 0 {
			insta.Spaceflight(c, *runSpaceflight)
		}

		if runTron.Seconds() > 0 {
			s := insta.NewScreen()
			tr := tron.NewGame(insta.ScreenWidth, insta.ScreenHeight)
			for {
				tr.Step(pads())
				tr.Paint(s)
				c.SetScreen(s)
				time.Sleep(1000 / time.Duration(*fps) * time.Millisecond)
			}
		}

		if runSnake.Seconds() > 0 {
			s := insta.NewScreen()
			sn := snake.NewGame(insta.ScreenWidth, insta.ScreenHeight)
			for {
				sn.Step(pads())
				sn.Paint(s)
				c.SetScreen(s)
				time.Sleep(1000 / time.Duration(*fps) * time.Millisecond)
			}
		}
		time.Sleep(20 * time.Millisecond)
	}
}
