package audio

import (
	"bufio"
	"fmt"
	"image/color"
	"io"
	"os"
	"strconv"
	"sync"

	"github.com/ktt-ol/go-insta"
)

type LevelGraph struct {
	input io.Reader

	step chan struct{}
	mu   sync.Mutex
	hist [insta.ScreenWidth]float64
}

func NewLevelGraph(r io.Reader) *LevelGraph {
	l := &LevelGraph{
		input: r,
		step:  make(chan struct{}, 1),
	}
	go l.Start()
	return l
}

func (g *LevelGraph) AddValue(v float64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	copy(g.hist[0:], g.hist[1:])
	g.hist[len(g.hist)-1] = v
}

func (g *LevelGraph) Start() {
	scanner := bufio.NewScanner(g.input)
	for scanner.Scan() {
		v, err := strconv.ParseFloat(scanner.Text(), 64)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error reading audio data:", err)
			continue
		}
		select {
		case g.step <- struct{}{}:
			g.AddValue(v)
		default:
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading audio data:", err)
	}
}

func (g *LevelGraph) Next() {
	<-g.step
}

func (g *LevelGraph) UpdateScreen(s *insta.Screen) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for x := 0; x < insta.ScreenWidth; x++ {
		for y := 0; y < insta.ScreenHeight; y++ {
			if (g.hist[x] / 30.0 * insta.ScreenHeight) > (insta.ScreenHeight - float64(y)) {
				r, g, b := insta.HsvToRgb(float64(y*3), 0.7, 0.8)
				s.Set(x, y, color.RGBA{uint8(r * 255), uint8(g * 255), uint8(b * 255), 128})
			} else {
				s.Set(x, y, color.RGBA{0, 0, 0, 128})
			}
		}
	}
}
