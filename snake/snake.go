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
	Head   Piece
	Tail   []Piece
	Length int
	Dir    Direction
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
	Player        Player
	tickDur       time.Duration
	idleTime      time.Time
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
		Player: Player{
			Head: Piece{
				Color: color.RGBA{255, 255, 0, 255},
				X:     insta.ScreenWidth / 2,
				Y:     insta.ScreenHeight / 2,
				Set:   true,
			},
			Dir:    dir,
			Length: 10,
		},
		tickDur:       time.Millisecond * 50,
		idleTime:      time.Now().Add(exitAfterIdle),
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

func (g *Game) Step(pads []insta.Pad) GameStatus {
	time.Sleep(g.tickDur)

	moved := false
	prev := g.Player.Dir
	if pads[0].Up() && prev != Down {
		g.Player.Dir = Up
		moved = true
	}
	if pads[0].Down() && prev != Up {
		g.Player.Dir = Down
		moved = true
	}
	if pads[0].Left() && prev != Right {
		g.Player.Dir = Left
		moved = true
	}
	if pads[0].Right() && prev != Left {
		g.Player.Dir = Right
		moved = true
	}
	if moved {
		g.idleTime = time.Now().Add(g.ExitAfterIdle)
	} else if g.idleTime.Before(time.Now()) {
		return Exit
	}

	if pads[0].Start() {
		return End
	}

	g.move()

	end := g.Player.PushHead()
	if end.Set {
		g.Field[end.Y][end.X].Snake = false
	}

	h := g.Player.Head
	if g.Field[h.Y][h.X].Snake {
		g.clear()
		g.Player.Dir = None
		return End
	} else if g.Field[h.Y][h.X].Fruit {
		g.Field[h.Y][h.X].Fruit = false
		g.Player.Head.Color = g.Field[h.Y][h.X].Color
		g.Player.Length += 10
		g.SpawnFruit()
		if g.tickDur > 20*time.Millisecond {
			g.tickDur -= 2 * time.Millisecond
		}
	}

	c := g.Field[h.Y][h.X]
	c.Snake = true
	c.Color = h.Color
	g.Field[h.Y][h.X] = c
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

	for _, t := range g.Player.Tail {
		img.Set(t.X, t.Y, t.Color)
	}
}

func (g *Game) PaintScore(img draw.Image) {
	fontWidth := 7
	fontHeight := 13

	text := fmt.Sprintf("%d", g.Player.Length)

	d := &font.Drawer{
		Dst:  img,
		Src:  image.White,
		Face: basicfont.Face7x13,
		Dot:  fixed.P(insta.ScreenWidth/2-fontWidth, insta.ScreenHeight/2+fontHeight/2),
	}
	d.DrawString(text)
}
