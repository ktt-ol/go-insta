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
	copy(pads[0][:], m.buf[1+6*0:6*0+6+1])
	copy(pads[1][:], m.buf[1+6*1:6*1+6+1])
	copy(pads[2][:], m.buf[1+6*2:6*2+6+1])
	copy(pads[3][:], m.buf[1+6*3:6*3+6+1])
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
	PSPAD_MASK_SELECT   = 0x01
	PSPAD_MASK_START    = 0x08
	PSPAD_MASK_UP       = 0x10
	PSPAD_MASK_RIGHT    = 0x20
	PSPAD_MASK_DOWN     = 0x40
	PSPAD_MASK_LEFT     = 0x80
	PSPAD_MASK_L2       = 0x01
	PSPAD_MASK_R2       = 0x02
	PSPAD_MASK_L1       = 0x04
	PSPAD_MASK_R1       = 0x08
	PSPAD_MASK_TRIANGLE = 0x10
	PSPAD_MASK_CIRCLE   = 0x20
	PSPAD_MASK_CROSS    = 0x40
	PSPAD_MASK_SQUARE   = 0x80
)

type Pad [6]byte

func (p Pad) Up() bool    { return (p[0] & PSPAD_MASK_UP) != 0 }
func (p Pad) Right() bool { return (p[0] & PSPAD_MASK_RIGHT) != 0 }
func (p Pad) Down() bool  { return (p[0] & PSPAD_MASK_DOWN) != 0 }
func (p Pad) Left() bool  { return (p[0] & PSPAD_MASK_LEFT) != 0 }

func (p Pad) Select() bool { return (p[0] & PSPAD_MASK_SELECT) != 0 }
func (p Pad) Start() bool  { return (p[0] & PSPAD_MASK_START) != 0 }

func (p Pad) R1() bool { return (p[1] & PSPAD_MASK_R1) != 0 }
func (p Pad) R2() bool { return (p[1] & PSPAD_MASK_R2) != 0 }
func (p Pad) L1() bool { return (p[1] & PSPAD_MASK_L1) != 0 }
func (p Pad) L2() bool { return (p[1] & PSPAD_MASK_L2) != 0 }

func (p Pad) Triangle() bool { return (p[1] & PSPAD_MASK_TRIANGLE) != 0 }
func (p Pad) Circle() bool   { return (p[1] & PSPAD_MASK_CIRCLE) != 0 }
func (p Pad) Cross() bool    { return (p[1] & PSPAD_MASK_CROSS) != 0 }
func (p Pad) Square() bool   { return (p[1] & PSPAD_MASK_SQUARE) != 0 }

func (p Pad) StickRight() (float32, float32) {
	return (float32(p[2]) - 128) / 128, (float32(p[3]) - 128) / 128
}
func (p Pad) StickLeft() (float32, float32) {
	return (float32(p[4]) - 128) / 128, (float32(p[5]) - 128) / 128
}

func (p *Pad) String() string {
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
