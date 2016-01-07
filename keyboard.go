package insta

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
)

func NewKeyboardPad() *KeyboardPad {
	kp := KeyboardPad{
		mu: &sync.Mutex{},
	}
	go kp.run()
	return &kp
}

type KeyboardPad struct {
	b  byte
	mu *sync.Mutex
}

func (m *KeyboardPad) run() {
	cmd := exec.Command("/bin/stty", "-f", "/dev/tty", "-icanon", "min", "1", "-echo")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			cmd := exec.Command("/bin/stty", "-f", "/dev/tty", "sane")
			if err := cmd.Run(); err != nil {
				log.Fatal(err)
			}
			os.Exit(1)
		}
	}()

	var keybuf [1]byte
	for {
		if n, err := os.Stdin.Read(keybuf[0:1]); err == nil && n == 1 {
			m.mu.Lock()
			m.b = keybuf[0]
			m.mu.Unlock()
		}
	}
}

func (m *KeyboardPad) Pads() []Pad {
	m.mu.Lock()
	defer m.mu.Unlock()
	ps := byteToPad(m.b)
	m.b = 0
	return ps
}

func byteToPad(b byte) []Pad {
	ps := make([]Pad, 4)

	switch b {
	case 'q':
		ps[0][1] = PsMaskSquare
	case 'w':
		ps[0][0] = PsMaskUp
	case 's':
		ps[0][0] = PsMaskDown
	case 'a':
		ps[0][0] = PsMaskLeft
	case 'd':
		ps[0][0] = PsMaskRight

	case 'u':
		ps[1][0] = PsMaskSquare
	case 'i':
		ps[1][0] = PsMaskUp
	case 'k':
		ps[1][0] = PsMaskDown
	case 'j':
		ps[1][0] = PsMaskLeft
	case 'l':
		ps[1][0] = PsMaskRight
	}

	return ps
}
