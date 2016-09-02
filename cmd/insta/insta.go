package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/tarm/serial"

	"github.com/ktt-ol/go-insta"
	"github.com/ktt-ol/go-insta/life"
	"github.com/ktt-ol/go-insta/snake"
	"github.com/ktt-ol/go-insta/tron"
)

var addrs = []string{
	// "192.168.3.101", // 1 at 0:f:17:10:53:d1
	// "192.168.3.102", // 2 at 0:f:17:10:53:ac
	// "192.168.3.103", // 3 at 0:f:17:10:53:b1
	// "192.168.3.104", // 4 at 0:f:17:10:53:b9
	// "192.168.3.105", // 5 at 0:f:17:10:53:90
	// "192.168.3.106", // 6 at 0:f:17:10:53:c3

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

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	var (
		fps            = flag.Int("fps", 25, "fps")
		runLife        = flag.Duration("life", 0, "run life for duration")
		runSnake       = flag.Duration("snake", 0, "run snake for duration")
		runTron        = flag.Duration("tron", 0, "run tron for duration")
		runSpaceflight = flag.Duration("spaceflight", 0, "run spaceflight for duration")
		runLogo        = flag.Bool("logo", false, "show mainframe logo")
		port           = flag.String("port", "", "serial port")
		joystick       = flag.Int("joystick", -1, "joystick id")
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
	} else if *joystick >= 0 {
		kp := insta.NewJoystick(*joystick)
		pads = kp.Pads
	} else {
		pads = func() []insta.Pad {
			return nil
		}
	}

	c.SetFPS(*fps)
	go c.Run()

	c.SetAfterglow(0.3)

	for {
		if *runLogo {
			c.SetAfterglow(0)
			insta.ShowImage(c, "img/mainframe-mod.png")
			c.SetAfterglow(0.4)
		}

		if runLife.Seconds() > 0 {
			l := life.NewLife(insta.ScreenWidth, insta.ScreenHeight)
			t := time.NewTicker(2 * time.Second)

			go func() {
				for _ = range t.C {
					l.AddRandomSpaceship()
				}
			}()

			s := insta.NewScreen()
			prev := s.Copy()
			sps := 10
			blendSteps := int(float32(*fps) / float32(sps))

			till := time.Now().Add(*runLife)

		lifeLoop:
			for time.Now().Before(till) {
				l.UpdateScreen(s)
				for _, img := range insta.BlendScreens(prev, s, blendSteps) {
					if pads()[0].Start() {
						break lifeLoop
					}
					c.SetScreen(img)
				}
				l.Step()
				prev, s = s, prev
			}

			t.Stop()
		}

		if runTron.Seconds() > 0 {
			s := insta.NewScreen()
			tr := tron.NewGame(insta.ScreenWidth, insta.ScreenHeight)
			for {
				tr.Step(pads())
				tr.Paint(s)
				c.SetScreenImmediate(s)
				time.Sleep(50 * time.Microsecond)
			}
		}

		if runSnake.Seconds() > 0 {
			c.SetAfterglow(0.2)
		SnakeLoop:
			for {
				s := insta.NewScreen()
				sn := snake.NewGame(insta.ScreenWidth, insta.ScreenHeight, *runSnake)
				for {
					status := sn.Step(pads())
					sn.Paint(s)
					c.SetScreenImmediate(s)
					if status == snake.End {
						time.Sleep(100 * time.Millisecond) // prevent screen from being dropped
						sn.PaintScore(s)
						c.SetScreenImmediate(s)
						time.Sleep(3 * time.Second)
						break
					}
					if status == snake.Exit {
						break SnakeLoop
					}
				}
			}
			c.SetAfterglow(0.3)
		}

		if runSpaceflight.Seconds() > 0 {
			insta.Spaceflight(c, *runSpaceflight)
		}
		time.Sleep(20 * time.Millisecond)
	}
}
