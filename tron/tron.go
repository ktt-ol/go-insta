package tron

import (
	"image/draw"

	"image/color"
)

type Pos struct {
	X, Y int
}

type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

type Player struct {
	Color color.RGBA
	Name  string
	Pos   Pos
	Dir   Direction
}

type Field struct {
	f    [][]*Player
	w, h int
}

func NewField(w, h int) *Field {
	f := make([][]*Player, h)
	for i := range f {
		f[i] = make([]*Player, w)
	}
	return &Field{f: f, w: w, h: h}
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
			&Player{Color: color.RGBA{255, 0, 0, 255}, Pos: Pos{X: 10, Y: 0}},
		},
	}
	return g
}

// movePos moves Pos in Dir, wraps around at field boundaries
func (f *Field) movePos(p Pos, d Direction) Pos {
	switch d {
	case Up:
		p.Y += 1
	case Down:
		p.Y -= 1
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

func (g *Game) Step() {
	for _, p := range g.Players {
		// if rand.Intn(4) == 0 {
		// 	p.Dir += Direction(rand.Intn(3) - 1)
		// 	p.Dir += 4
		// 	p.Dir %= 4
		// }
		p.Pos = g.Field.movePos(p.Pos, p.Dir)
	}
	for _, p := range g.Players {
		g.Field.f[p.Pos.Y][p.Pos.X] = p
	}
}

func (g *Game) Paint(img draw.Image) {
	for y := 0; y < g.Field.h; y++ {
		for x := 0; x < g.Field.w; x++ {
			if g.Field.f[y][x] != nil {
				img.Set(x, y, g.Field.f[y][x].Color)
			}
		}
	}
}
