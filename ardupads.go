package insta

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/tarm/serial"
)

func main() {
	port := flag.String("port", "", "serial port")
	flag.Parse()

	cfg := &serial.Config{}
	cfg.Baud = 57600
	cfg.Name = *port

	s, err := serial.OpenPort(cfg)
	if err != nil {
		log.Fatal(err)
	}

	mp := NewMultiPad(s)
	t := time.Tick(100 * time.Millisecond)
	for _ = range t {
		pads := mp.Pads()
		fmt.Println(0, &pads[0], 1, &pads[1], 2, &pads[2], 3, &pads[3])
	}
}

func NewMultiPad(s *serial.Port) *MultiPad {
	p := MultiPad{
		s:   s,
		buf: make([]byte, 1+6*4+2),
		mu:  &sync.Mutex{},
	}
	if s != nil {
		go p.run()
	}
	return &p
}

type MultiPad struct {
	s   *serial.Port
	buf []byte
	mu  *sync.Mutex
}

func (m *MultiPad) run() {
	b := make([]byte, 1+6*4+2)
	for {
		readMessage(m.s, b)
		if b[0] != 0xaa || b[25] != 0xaa || b[26] != '\n' {
			seekNewline(m.s)
			continue
		}
		m.mu.Lock()
		copy(m.buf[:], b[:])
		m.mu.Unlock()
	}
}

func (m *MultiPad) Pads() []Pad {
	pads := make([]Pad, 4)
	m.mu.Lock()
	defer m.mu.Unlock()
	for i := range pads {
		pad := PSController{}
		copy(pad[:], m.buf[1+6*i:6*i+6+1])
		pads[0] = &pad
	}
	return pads
}

func readMessage(s *serial.Port, b []byte) {
	var i = 0
	for {
		n, err := s.Read(b[i:])
		if err != nil {
			log.Fatal(err)
		}
		i += n
		if len(b) == i {
			break
		}
	}
}

func seekNewline(s *serial.Port) error {
	b := []byte{0}
	for {
		if _, err := s.Read(b); err != nil {
			return err
		}
		if b[0] == 0xaa {
			if _, err := s.Read(b); err != nil {
				return err
			}
			if b[0] == '\n' {
				break
			}
		}
	}
	return nil
}

const (
	PsMaskSelect   = 0x01
	PsMaskStart    = 0x08
	PsMaskUp       = 0x10
	PsMaskRight    = 0x20
	PsMaskDown     = 0x40
	PsMaskLeft     = 0x80
	PsMaskL2       = 0x01
	PsMaskR2       = 0x02
	PsMaskL1       = 0x04
	PsMaskR1       = 0x08
	PsMaskTriangle = 0x10
	PsMaskCircle   = 0x20
	PsMaskCross    = 0x40
	PsMaskSquare   = 0x80
)

type Pad interface {
	Up() bool
	Down() bool
	Left() bool
	Right() bool

	Start() bool
	Select() bool

	South() bool
	East() bool
	North() bool
	West() bool
}

type PSController [6]byte

func (p *PSController) Up() bool    { return (p[0] & PsMaskUp) != 0 }
func (p *PSController) Right() bool { return (p[0] & PsMaskRight) != 0 }
func (p *PSController) Down() bool  { return (p[0] & PsMaskDown) != 0 }
func (p *PSController) Left() bool  { return (p[0] & PsMaskLeft) != 0 }

func (p *PSController) Select() bool { return (p[0] & PsMaskSelect) != 0 }
func (p *PSController) Start() bool  { return (p[0] & PsMaskStart) != 0 }

func (p *PSController) R1() bool { return (p[1] & PsMaskR1) != 0 }
func (p *PSController) R2() bool { return (p[1] & PsMaskR2) != 0 }
func (p *PSController) L1() bool { return (p[1] & PsMaskL1) != 0 }
func (p *PSController) L2() bool { return (p[1] & PsMaskL2) != 0 }

func (p *PSController) Triangle() bool { return (p[1] & PsMaskTriangle) != 0 }
func (p *PSController) Circle() bool   { return (p[1] & PsMaskCircle) != 0 }
func (p *PSController) Cross() bool    { return (p[1] & PsMaskCross) != 0 }
func (p *PSController) Square() bool   { return (p[1] & PsMaskSquare) != 0 }

func (p *PSController) North() bool { return (p[1] & PsMaskTriangle) != 0 }
func (p *PSController) East() bool  { return (p[1] & PsMaskCircle) != 0 }
func (p *PSController) South() bool { return (p[1] & PsMaskCross) != 0 }
func (p *PSController) West() bool  { return (p[1] & PsMaskSquare) != 0 }

func (p *PSController) StickRight() (float32, float32) {
	return (float32(p[2]) - 128) / 128, (float32(p[3]) - 128) / 128
}
func (p *PSController) StickLeft() (float32, float32) {
	return (float32(p[4]) - 128) / 128, (float32(p[5]) - 128) / 128
}

func (p *PSController) String() string {
	b := bytes.Buffer{}
	if p.Up() {
		b.WriteString("u ")
	}
	if p.Right() {
		b.WriteString("r ")
	}
	if p.Down() {
		b.WriteString("d ")
	}
	if p.Left() {
		b.WriteString("l ")
	}

	if p.Triangle() {
		b.WriteString("t ")
	}
	if p.Circle() {
		b.WriteString("c ")
	}
	if p.Cross() {
		b.WriteString("x ")
	}
	if p.Square() {
		b.WriteString("q ")
	}

	if lr, ud := p.StickLeft(); lr != 0.0 || ud != 0.0 {
		fmt.Fprintf(&b, "L %.2f %.2f ", lr, ud)
	}
	if lr, ud := p.StickRight(); lr != 0.0 || ud != 0.0 {
		fmt.Fprintf(&b, "R %.2f %.2f ", lr, ud)
	}

	return string(b.Bytes())
}

type NullPad struct{}

func (p *NullPad) Up() bool    { return false }
func (p *NullPad) Right() bool { return false }
func (p *NullPad) Down() bool  { return false }
func (p *NullPad) Left() bool  { return false }

func (p *NullPad) Select() bool { return false }
func (p *NullPad) Start() bool  { return false }

func (p *NullPad) North() bool { return false }
func (p *NullPad) East() bool  { return false }
func (p *NullPad) South() bool { return false }
func (p *NullPad) West() bool  { return false }
