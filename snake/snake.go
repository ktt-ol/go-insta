package snake

import (
	"image/draw"

	"github.com/ktt-ol/go-insta"

	"image/color"
)

type Piece struct {
	X, Y  int
	Color color.RGBA
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
	Head Piece
	Tail []Piece
	Dir  Direction
}

type Cell struct {
	Color color.RGBA
	Fruit bool
	Snake bool
}

type Field [][]Cell

type Game struct {
	Width, Height int
	Field         Field
	Player        Player
}

func NewGame(w, h int) *Game {
	f := make([][]Cell, h)
	for i := range f {
		f[i] = make([]Cell, w)
	}
	g := &Game{
		Field:  f,
		Width:  w,
		Height: h,
		Player: Player{
			Head: Piece{
				Color: color.RGBA{255, 255, 0, 255},
				X:     insta.ScreenWidth / 2,
				Y:     insta.ScreenHeight / 2,
			},
			Dir:  Down,
			Tail: make([]Piece, 10),
		},
	}
	return g
}

func (g *Game) clear() {
	g.Field = make([][]Cell, g.Height)
	for i := range g.Field {
		g.Field[i] = make([]Cell, g.Width)
	}
}

func (g *Game) move() {
	switch g.Player.Dir {
	case Up:
		g.Player.Head.Y -= 1
	case Down:
		g.Player.Head.Y += 1
	case Right:
		g.Player.Head.X += 1
	case Left:
		g.Player.Head.X -= 1
	}
	g.Player.Head.X += g.Width
	g.Player.Head.X %= g.Width
	g.Player.Head.Y += g.Height
	g.Player.Head.Y %= g.Height
}

func (g *Game) Step(pads []insta.Pad) {
	if pads[0].Up() || pads[0].Triangle() {
		g.Player.Dir = Up
	}
	if pads[0].Down() || pads[0].Cross() {
		g.Player.Dir = Down
	}
	if pads[0].Left() || pads[0].Square() {
		g.Player.Dir = Left
	}
	if pads[0].Right() || pads[0].Circle() {
		g.Player.Dir = Right
	}

	g.move()
	// if g.Field.f[g.Player.Head.Y][g.Player.Head.X].Set {
	// 	g.Field.Clear()
	// 	for _, p := range g.Players {
	// 		g.Player.Dir = None
	// 	}
	// 	return
	// }

	c := g.Field[g.Player.Head.Y][g.Player.Head.X]
	c.Snake = true
	c.Color = g.Player.Head.Color
	g.Field[g.Player.Head.Y][g.Player.Head.X] = c
}

func (g *Game) Paint(img draw.Image) {
	for y := range g.Field {
		for x := range g.Field[y] {
			if g.Field[y][x].Snake {
				img.Set(x, y, g.Field[y][x].Color)
			} else {
				img.Set(x, y, color.RGBA{0, 0, 0, 0})
			}
		}
	}
}
