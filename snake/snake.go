package snake

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math/rand"
	"time"

	"golang.org/x/image/font/basicfont"

	"github.com/ktt-ol/go-insta"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
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
	Head     Piece
	Tail     []Piece
	Length   int
	Score    int
	Dir      Direction
	idleTime time.Time
}

func (p *Player) PushHead() Piece {
	var end Piece
	if len(p.Tail) < p.Length {
		p.Tail = append(p.Tail, Piece{})
		copy(p.Tail[1:], p.Tail)
		p.Tail[0] = p.Head
	} else {
		end = p.Tail[len(p.Tail)-1]
		for i := len(p.Tail) - 2; i >= 0; i-- {
			p.Tail[i+1].X = p.Tail[i].X
			p.Tail[i+1].Y = p.Tail[i].Y
		}
		p.Tail[0] = p.Head
	}
	return end
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
	Players       []*Player
	tickDur       time.Duration
	ExitAfterIdle time.Duration
}

type GameStatus int

const (
	Running GameStatus = iota
	End
	Exit
)

func NewGame(w, h int, exitAfterIdle time.Duration) *Game {
	f := make([][]Cell, h)
	for i := range f {
		f[i] = make([]Cell, w)
	}

	dir := Down
	switch rand.Intn(4) {
	case 0:
		dir = Up
	case 1:
		dir = Left
	case 2:
		dir = Right
	}
	g := &Game{
		Field:  f,
		Width:  w,
		Height: h,
		Players: []*Player{
			&Player{
				Head: Piece{
					Color: color.RGBA{255, 255, 0, 255},
					X:     insta.ScreenWidth / 2,
					Y:     insta.ScreenHeight / 4,
					Set:   true,
				},
				Dir:      dir,
				Length:   10,
				idleTime: time.Now().Add(exitAfterIdle),
			},
			&Player{
				Head: Piece{
					Color: color.RGBA{255, 0, 255, 255},
					X:     insta.ScreenWidth / 2,
					Y:     insta.ScreenHeight - insta.ScreenHeight/4,
					Set:   true,
				},
				Dir:      dir,
				Length:   10,
				idleTime: time.Now().Add(exitAfterIdle),
			},
		},
		tickDur:       time.Millisecond * 50,
		ExitAfterIdle: exitAfterIdle,
	}

	g.SpawnFruit()
	g.SpawnFruit()
	g.SpawnFruit()

	return g
}

func (g *Game) clear() {
	g.Field = make([][]Cell, g.Height)
	for i := range g.Field {
		g.Field[i] = make([]Cell, g.Width)
	}
}

func (g *Game) move(p *Player) {
	switch p.Dir {
	case Up:
		p.Head.Y -= 1
	case Down:
		p.Head.Y += 1
	case Right:
		p.Head.X += 1
	case Left:
		p.Head.X -= 1
	}
	p.Head.X += g.Width
	p.Head.X %= g.Width
	p.Head.Y += g.Height
	p.Head.Y %= g.Height
}

func (g *Game) Step(pads []insta.Pad) GameStatus {
	time.Sleep(g.tickDur)

	activePlayer := false
	for i, p := range g.Players {
		if p.idleTime.IsZero() {
			continue
		}
		activePlayer = true

		moved := false

		prev := p.Dir
		if pads[i].Up() && prev != Down {
			p.Dir = Up
			moved = true
		}
		if pads[i].Down() && prev != Up {
			p.Dir = Down
			moved = true
		}
		if pads[i].Left() && prev != Right {
			p.Dir = Left
			moved = true
		}
		if pads[i].Right() && prev != Left {
			p.Dir = Right
			moved = true
		}
		if moved {
			p.idleTime = time.Now().Add(g.ExitAfterIdle)
		} else if p.idleTime.Before(time.Now()) {
			p.idleTime = time.Time{}
			for _, p := range p.Tail {
				g.Field[p.Y][p.X].Snake = false
			}
			g.Field[p.Head.Y][p.Head.X].Snake = false
			p.Tail = []Piece{}
			g.Players[i] = p
			continue
		}

		if pads[0].Start() {
			return End
		}

		g.move(p)

		end := p.PushHead()
		if end.Set {
			g.Field[end.Y][end.X].Snake = false
		}
	}

	// check for head crash, both players loose
	for i, p := range g.Players {
		if p.idleTime.IsZero() {
			continue
		}
		for _, o := range g.Players[i+1:] {
			if o.idleTime.IsZero() {
				continue
			}
			if p.Head.X == o.Head.X && p.Head.Y == o.Head.Y {
				p.Score -= 50
				o.Score -= 50
				g.clear()
				return End
			}
		}
	}

	for _, p := range g.Players {
		if p.idleTime.IsZero() {
			continue
		}
		h := p.Head
		if g.Field[h.Y][h.X].Snake {
			g.clear()
			p.Dir = None
			p.Score -= 50
			return End
		} else if g.Field[h.Y][h.X].Fruit {
			g.Field[h.Y][h.X].Fruit = false
			p.Head.Color = g.Field[h.Y][h.X].Color
			p.Length += 10
			p.Score += 10
			g.SpawnFruit()
			if g.tickDur > 20*time.Millisecond {
				g.tickDur -= 1 * time.Millisecond
			}
		}

		c := g.Field[h.Y][h.X]
		c.Snake = true
		c.Color = h.Color
		g.Field[h.Y][h.X] = c

	}
	if !activePlayer {
		return Exit
	}
	return Running
}

func (g *Game) SpawnFruit() {
	for {
		y := rand.Intn(len(g.Field))
		x := rand.Intn(len(g.Field[0]))
		if !g.Field[y][x].Snake && !g.Field[y][x].Fruit {
			g.Field[y][x].Fruit = true
			g.Field[y][x].Color = insta.HsvToColor(rand.Float64()*360, 0.8, 0.5)
			return
		}
	}
}

func (g *Game) Paint(img draw.Image) {
	for y := range g.Field {
		for x := range g.Field[y] {
			if g.Field[y][x].Fruit {
				img.Set(x, y, g.Field[y][x].Color)
			} else {
				img.Set(x, y, color.RGBA{0, 0, 0, 0})
			}
		}
	}

	for _, p := range g.Players {
		for _, t := range p.Tail {
			img.Set(t.X, t.Y, t.Color)
		}
	}
}

func (g *Game) PaintScore(img draw.Image) {
	fontWidth := 7
	fontHeight := 13

	scores := " "
	for _, p := range g.Players {
		if p.Score == 0 {
			scores += "    "
		} else {
			scores += fmt.Sprintf("%3d ", p.Score)
		}
	}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.White,
		Face: basicfont.Face7x13,
		Dot:  fixed.P(insta.ScreenWidth/2-(fontWidth*len(scores)/2), insta.ScreenHeight/2+fontHeight/2),
	}
	d.DrawString(scores)
}
