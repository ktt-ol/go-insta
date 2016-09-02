package insta

import (
	"log"

	"github.com/simulatedsimian/joystick"
)

type Joystick struct {
	sticks []joystick.Joystick
}

func NewJoystick(ids []int) *Joystick {
	js := &Joystick{}
	for _, id := range ids {
		j, err := joystick.Open(id)
		if err != nil {
			log.Fatalf("unable to open joystick %d: %s", id, err)
		}
		js.sticks = append(js.sticks, j)
	}
	return js
}

func (js *Joystick) Pads() []Pad {
	pads := make([]Pad, 4)
	for i, j := range js.sticks {
		s, _ := j.Read()
		p := JoystickState{s}
		pads[i] = &p
	}
	return pads
}

type JoystickState struct {
	joystick.State
}

func (js *JoystickState) Up() bool    { return js.AxisData[1] < -1000 }
func (js *JoystickState) Right() bool { return js.AxisData[0] > 1000 }
func (js *JoystickState) Down() bool  { return js.AxisData[1] > 1000 }
func (js *JoystickState) Left() bool  { return js.AxisData[0] < -1000 }

func (js *JoystickState) Select() bool { return (js.Buttons & 256) != 0 }
func (js *JoystickState) Start() bool  { return (js.Buttons & 512) != 0 }

func (js *JoystickState) North() bool { return (js.Buttons & 1) != 0 }
func (js *JoystickState) East() bool  { return (js.Buttons & 2) != 0 }
func (js *JoystickState) South() bool { return (js.Buttons & 4) != 0 }
func (js *JoystickState) West() bool  { return (js.Buttons & 8) != 0 }

func (js *JoystickState) LPad() bool { return (js.Buttons & 16) != 0 }
func (js *JoystickState) RPad() bool { return (js.Buttons & 32) != 0 }
