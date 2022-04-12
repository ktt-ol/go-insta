package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"strings"
	"time"

	"github.com/ktt-ol/go-insta/srv"

	"github.com/tarm/serial"

	"github.com/ktt-ol/go-insta"
	"github.com/ktt-ol/go-insta/audio"
	"github.com/ktt-ol/go-insta/life"
	"github.com/ktt-ol/go-insta/snake"
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
		runSpaceflight = flag.Duration("spaceflight", 0, "run spaceflight for duration")
		runLogo        = flag.Bool("logo", false, "show mainframe logo")
		runAudio       = flag.Duration("audio", 0, "audio graph duration")
		runRainbow     = flag.Duration("rainbow", 0, "rainbow duration")
		runGifs        = flag.Duration("gifs", 0, "gif repeat duration")
		runServer      = flag.Bool("server", false, "start TCP server on port 2323, accepting images")
		audioDevice    = flag.String("audiodevice", "", "serial port of audio device")
		port           = flag.String("port", "", "serial port")
		joystick       = flag.String("joystick", "", "joystick ids")
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
	} else if *joystick != "" {
		var ids []int
		for _, idStr := range strings.Split(*joystick, ",") {
			id, err := strconv.ParseInt(idStr, 10, 32)
			if err != nil {
				log.Fatalf("invalid -joystick option", err)
			}
			ids = append(ids, int(id))
		}
		kp := insta.NewJoystick(ids)
		pads = kp.Pads
	} else {
		pads = func() []insta.Pad {
			return []insta.Pad{
				&insta.NullPad{},
				&insta.NullPad{},
			}
		}
	}

	var audioGraph *audio.LevelGraph
	if runAudio.Seconds() > 0 && *audioDevice != "" {
		cfg := &serial.Config{}
		cfg.Baud = 115200
		cfg.Name = *audioDevice
		cfg.ReadTimeout = time.Second * 5
		var err error
		audioInput, err := serial.OpenPort(cfg)
		if err != nil {
			log.Println("warn: skipping audio", err)
		} else {
			audioGraph = audio.NewLevelGraph(audioInput)
		}
	}

	c.SetFPS(*fps)
	go c.Run()

	c.SetAfterglow(0.3)

	if *runServer {
		srv.Server(c)
		return
	}

	for {
		if runRainbow.Seconds() > 0 {
			insta.Rainbow(c, *runRainbow)
		}

		if *runLogo {
			c.SetAfterglow(0)
			insta.ScrollImage(c, "img/mainframe-mod.png")
			c.SetAfterglow(0.4)
		}

		if runGifs.Seconds() > 0 {
			insta.RandomGif(c, "gifs", *runGifs)
			time.Sleep(100 * time.Millisecond)
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

		if runGifs.Seconds() > 0 {
			insta.RandomGif(c, "gifs", *runGifs)
			time.Sleep(200 * time.Millisecond)
		}

		if runAudio.Seconds() > 0 && audioGraph != nil {
			till := time.Now().Add(*runAudio)
			s := insta.NewScreen()

			c.SetAfterglow(0.1)
			for time.Now().Before(till) {
				audioGraph.Next()
				audioGraph.UpdateScreen(s)
				c.SetScreenImmediate(s)
				time.Sleep(50 * time.Millisecond)
			}
		}

		if runGifs.Seconds() > 0 {
			insta.RandomGif(c, "gifs", *runGifs)
			time.Sleep(200 * time.Millisecond)
		}

		if runSnake.Seconds() > 0 {
			c.SetAfterglow(0.2)
			sn := snake.NewGame(insta.ScreenWidth, insta.ScreenHeight, *runSnake)
		SnakeLoop:
			for {
				s := insta.NewScreen()
				sn.Init()
				for {
					status := sn.Step(pads())
					sn.Paint(s)
					c.SetScreenImmediate(s)
					if status == snake.End {
						time.Sleep(250 * time.Millisecond)
						// wait to prevent score screen from being overdrawn
						// by last game screen
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
