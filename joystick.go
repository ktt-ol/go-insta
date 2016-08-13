package insta

import (
	"log"

	"github.com/simulatedsimian/joystick"
)

type Joystick struct {
	stick joystick.Joystick
}

func NewJoystick(id int) *Joystick {
	j, err := joystick.Open(id)
	if err != nil {
		log.Fatalf("unable to open joystick %d: %s", id, err)
	}

	return &Joystick{stick: j}
}

func (j *Joystick) Pads() []Pad {
	pads := make([]Pad, 4)
	s, _ := j.stick.Read()

	p := Pad{}
	if s.AxisData[0] < -1000 {
		p[0] |= PsMaskLeft
	} else if s.AxisData[0] > 1000 {
		p[0] |= PsMaskRight
	}

	if s.AxisData[1] < -1000 {
		p[0] |= PsMaskUp
	} else if s.AxisData[1] > 1000 {
		p[0] |= PsMaskDown
	}
	pads[0] = p
	return pads
}
