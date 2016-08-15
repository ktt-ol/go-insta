package tron

import (
	"image/draw"

	"github.com/ktt-ol/go-insta"

	"image/color"
)

type Pos struct {
	X, Y int
}

type Direction int

const (
	None Direction = iota
	Up
	Right
	Down
	Left
)

func (d Direction) RotateLeft() Direction {
	if d == Up {
		return Left
	}
	return Direction(int(d) - 1)
}

func (d Direction) RotateRight() Direction {
	if d == Left {
		return Up
	}
	return Direction(int(d) + 1)
}

type Player struct {
	Color color.RGBA
	Name  string
	Pos   Pos
	Dir   Direction
}

type Cell struct {
	Color     color.RGBA
	Set       bool
	PlayerIdx int
}

type Field struct {
	f    [][]Cell
	w, h int
}

func NewField(w, h int) *Field {
	f := make([][]Cell, h)
	for i := range f {
		f[i] = make([]Cell, w)
	}
	return &Field{f: f, w: w, h: h}
}

func (f *Field) Clear() {
	f.f = make([][]Cell, f.h)
	for i := range f.f {
		f.f[i] = make([]Cell, f.w)
	}
}

type Game struct {
	Field   *Field
	Players []*Player
}

func NewGame(w, h int) *Game {
	g := &Game{
		Field: NewField(w, h),
		Players: []*Player{
			&Player{Color: color.RGBA{255, 255, 0, 255}, Pos: Pos{X: 0, Y: 0}},
			&Player{Color: color.RGBA{255, 0, 255, 255}, Pos: Pos{X: insta.ScreenWidth - 1, Y: 0}},
			&Player{Color: color.RGBA{0, 255, 255, 55}, Pos: Pos{X: 0, Y: insta.ScreenHeight - 1}},
			&Player{Color: color.RGBA{255, 0, 0, 255}, Pos: Pos{X: insta.ScreenWidth - 1, Y: insta.ScreenHeight - 1}},
		},
	}
	return g
}

// movePos moves Pos in Dir, wraps around at field boundaries
func (f *Field) movePos(p Pos, d Direction) Pos {
	switch d {
	case Up:
		p.Y -= 1
	case Down:
		p.Y += 1
	case Right:
		p.X += 1
	case Left:
		p.X -= 1
	}
	p.X += f.w
	p.X %= f.w
	p.Y += f.h
	p.Y %= f.h
	return p
}

func (g *Game) Step(pads []insta.Pad) {
	for i, p := range g.Players {
		if pads[i].Up() {
			p.Dir = Up
		}
		if pads[i].Down() {
			p.Dir = Down
		}
		if pads[i].Left() {
			p.Dir = Left
		}
		if pads[i].Right() {
			p.Dir = Right
		}

		if p.Dir == None {
			continue
		}

		p.Pos = g.Field.movePos(p.Pos, p.Dir)
		if g.Field.f[p.Pos.Y][p.Pos.X].Set {
			g.Field.Clear()
			for _, p := range g.Players {
				p.Dir = None
			}
			return
		}
	}

	for i, p := range g.Players {
		c := g.Field.f[p.Pos.Y][p.Pos.X]
		c.Set = true
		c.PlayerIdx = i
		c.Color = p.Color
		g.Field.f[p.Pos.Y][p.Pos.X] = c
	}
}

func (g *Game) Paint(img draw.Image) {
	for y := 0; y < g.Field.h; y++ {
		for x := 0; x < g.Field.w; x++ {
			if g.Field.f[y][x].Set {
				img.Set(x, y, g.Field.f[y][x].Color)
			} else {
				img.Set(x, y, color.RGBA{0, 0, 0, 0})
			}
		}
	}
}
