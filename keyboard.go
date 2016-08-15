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
	pad0 := PSController{}
	pad1 := PSController{}

	switch b {
	case 'q':
		pad0[1] = PsMaskSquare
	case 'w':
		pad0[0] = PsMaskUp
	case 's':
		pad0[0] = PsMaskDown
	case 'a':
		pad0[0] = PsMaskLeft
	case 'd':
		pad0[0] = PsMaskRight

	case 'u':
		pad1[0] = PsMaskSquare
	case 'i':
		pad1[0] = PsMaskUp
	case 'k':
		pad1[0] = PsMaskDown
	case 'j':
		pad1[0] = PsMaskLeft
	case 'l':
		pad1[0] = PsMaskRight
	}

	ps := make([]Pad, 4)
	ps[0] = &pad0
	ps[1] = &pad1
	return ps
}
