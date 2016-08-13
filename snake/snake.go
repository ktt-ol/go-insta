package snake

import (
	"image/draw"

	"github.com/ktt-ol/go-insta"

	"image/color"
)

type Piece struct {
	X, Y  int
	Set   bool
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
	f[10][20].Fruit = true
	f[20][10].Fruit = true
	f[30][30].Fruit = true
	g := &Game{
		Field:  f,
		Width:  w,
		Height: h,
		Player: Player{
			Head: Piece{
				Color: color.RGBA{255, 255, 0, 255},
				X:     insta.ScreenWidth / 2,
				Y:     insta.ScreenHeight / 2,
				Set:   true,
			},
			Dir:  Down,
			Tail: make([]Piece, 20),
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

	p := g.Player.Tail[len(g.Player.Tail)-1]
	if p.Set {
		g.Field[p.Y][p.X].Snake = false
	}
	copy(g.Player.Tail[1:], g.Player.Tail)
	g.Player.Tail[0] = g.Player.Head

	h := g.Player.Head
	if g.Field[h.Y][h.X].Snake {
		g.clear()
		g.Player.Dir = None
		return
	} else if g.Field[h.Y][h.X].Fruit {
		g.Field[h.Y][h.X].Fruit = false
		oldtail := g.Player.Tail
		g.Player.Tail = make([]Piece, len(oldtail)+10)
		copy(g.Player.Tail, oldtail)
	}

	c := g.Field[h.Y][h.X]
	c.Snake = true
	c.Color = h.Color
	g.Field[h.Y][h.X] = c
}

func (g *Game) Paint(img draw.Image) {
	for y := range g.Field {
		for x := range g.Field[y] {
			if g.Field[y][x].Snake {
				img.Set(x, y, g.Field[y][x].Color)
			} else if g.Field[y][x].Fruit {
				img.Set(x, y, color.White)
			} else {
				img.Set(x, y, color.RGBA{0, 0, 0, 0})
			}
		}
	}
}
