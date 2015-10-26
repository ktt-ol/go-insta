// An implementation of Conway's Game of Life.
// Adapted from Go (golang.org/doc/play/life.go) to calculate
// colored cells.
// Copyright (c) 2012 The Go Authors. All rights reserved.
package life

import (
	"bytes"
	"math"
	"math/rand"
)

type Cell struct {
	Hue   float32
	Count int
	Alive bool
}

// Field represents a two-dimensional field of cells.
type Field struct {
	s    [][]Cell
	w, h int
}

// NewField returns an empty field of the specified width and height.
func NewField(w, h int) *Field {
	s := make([][]Cell, h)
	for i := range s {
		s[i] = make([]Cell, w)
	}
	return &Field{s: s, w: w, h: h}
}

// Set sets the state of the specified cell to the given value.
func (f *Field) Set(x, y int, b Cell) {
	f.s[y][x] = b
}

// Alive reports whether the specified cell is alive.
func (f *Field) Alive(x, y int) bool {
	return f.Cell(x, y).Alive
}

// Cell returns the specified cell.
// If the x or y coordinates are outside the field boundaries they are wrapped
// toroidally. For instance, an x value of -1 is treated as width-1.
func (f *Field) Cell(x, y int) Cell {
	x += f.w
	x %= f.w
	y += f.h
	y %= f.h
	return f.s[y][x]
}

// Next returns the state of the specified cell at the next time step.
func (f *Field) Next(x, y int) Cell {
	// Count the adjacent cells that are alive.
	alive := 0
	// Add all hue values
	hue := float32(0.0)
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if (j != 0 || i != 0) && f.Alive(x+i, y+j) {
				hue += f.Cell(x+i, y+j).Hue
				alive++
			}
		}
	}
	// Return next state according to the game rules:
	//   exactly 3 neighbors: alive,
	//   exactly 2 neighbors: maintain current state,
	//   otherwise: off.
	if alive == 3 || alive == 2 && f.Alive(x, y) {
		c := f.Cell(x, y)
		c.Count = alive
		if c.Alive {
			alive += 1
			hue += c.Hue
		}
		c.Alive = true
		// Calculate median Hue of all alive Cells, shift by fixed value
		// fmt.Println(hue, alive, c.Hue, float32(math.Mod(float64(hue)/float64(alive)+2.0, 360)))
		c.Hue = float32(math.Mod(float64(hue)/float64(alive)+2.0, 360))
		return c
	}
	return Cell{}
}

// Life stores the state of a round of Conway's Game of Life.
type Life struct {
	a, b *Field
	w, h int
}

// NewLife returns a new Life game state with a random initial state.
func NewLife(w, h int) *Life {
	a := NewField(w, h)
	for i := 0; i < (w * h / 4); i++ {
		a.Set(rand.Intn(w), rand.Intn(h), Cell{Hue: rand.Float32() * 360, Count: 1, Alive: true})
	}
	// a.Set(10, 10, Cell{Alive: true, Count: 2, Hue: 100})
	// a.Set(10, 11, Cell{Alive: true, Count: 3, Hue: 100})
	// a.Set(10, 12, Cell{Alive: true, Count: 2, Hue: 100})
	return &Life{
		a: a, b: NewField(w, h),
		w: w, h: h,
	}
}

// Step advances the game by one instant, recomputing and updating all cells.
func (l *Life) Step() {
	// Update the state of the next field (b) from the current field (a).
	for y := 0; y < l.h; y++ {
		for x := 0; x < l.w; x++ {
			l.b.Set(x, y, l.a.Next(x, y))
		}
	}
	// Swap fields a and b.
	l.a, l.b = l.b, l.a
}

func (l *Life) Field() *Field {
	return l.a
}

// String returns the game board as a string.
func (l *Life) String() string {
	var buf bytes.Buffer
	for y := 0; y < l.h; y++ {
		for x := 0; x < l.w; x++ {
			b := byte(' ')
			if l.a.Alive(x, y) {
				b = '*'
			}
			buf.WriteByte(b)
		}
		buf.WriteByte('\n')
	}
	return buf.String()
}
