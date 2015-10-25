package insta

import (
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"strconv"

	"github.com/lucasb-eyer/go-colorful"

	"github.com/ktt-ol/go-insta/life"
)

const (
	PANEL_WIDTH  = 18
	PANEL_HWIDTH = PANEL_WIDTH / 2
	PANEL_HEIGHT = 18

	PANELS_X = 4
	PANELS_Y = 3

	SCREEN_WIDTH  = PANEL_WIDTH * PANELS_X
	SCREEN_HEIGHT = PANEL_HEIGHT * PANELS_Y

	PIXEL_STRIDE   = 3
	PANEL_STRIDE   = PIXEL_STRIDE * PANEL_WIDTH
	LINE_STRIDE    = PANEL_STRIDE * PANELS_X
	PANEL_Y_STRIDE = LINE_STRIDE * PANEL_HEIGHT
)

type Screen struct {
	Pix []uint8
}

func NewScreen() *Screen {
	return &Screen{
		Pix: make([]uint8, SCREEN_WIDTH*SCREEN_HEIGHT*PIXEL_STRIDE),
	}
}

func (s *Screen) Set(x, y int, c color.RGBA) {
	offset := y*LINE_STRIDE + x*PIXEL_STRIDE
	s.Pix[offset] = c.R
	s.Pix[offset+1] = c.G
	s.Pix[offset+2] = c.B
}

func (s *Screen) At(x, y int) color.Color {
	offset := y*LINE_STRIDE + x*PIXEL_STRIDE
	return color.RGBA{
		s.Pix[offset],
		s.Pix[offset+1],
		s.Pix[offset+2],
		255,
	}
}

func (s *Screen) Bounds() image.Rectangle {
	return image.Rect(0, 0, SCREEN_WIDTH, SCREEN_HEIGHT)
}

func (s *Screen) ColorModel() color.Model {
	return color.RGBAModel
}

func (s *Screen) String() string {
	var b []byte
	offset := 0
	for y := 0; y < SCREEN_HEIGHT; y++ {
		for x := 0; x < SCREEN_WIDTH; x++ {
			for c := 0; c < PIXEL_STRIDE; c++ {
				v := s.Pix[offset]
				if v < 16 {
					b = append(b, '0')
				}
				b = strconv.AppendInt(b, int64(v), 16)
				offset += 1
			}
			b = append(b, ' ')
		}
		b = append(b, '\n')
	}
	return string(b)
}

// Panel returns the left and light part of panel x/y.
func (s *Screen) Panel(x, y int) ([486]uint8, [486]uint8) {
	l := [486]uint8{}
	r := [486]uint8{}

	panelOffset := x*PANEL_STRIDE + y*PANEL_Y_STRIDE
	for py := 0; py < PANEL_HEIGHT; py++ {
		offset := panelOffset + py*LINE_STRIDE
		copy(l[py*PANEL_WIDTH/2:], s.Pix[offset:offset+PANEL_HWIDTH*PIXEL_STRIDE])
		copy(r[py*PANEL_WIDTH/2:], s.Pix[offset+PANEL_HWIDTH:offset+PANEL_WIDTH*PIXEL_STRIDE])

	}
	return l, r
}

func LifeToScreen(l *life.Life, s *Screen) {
	f := l.Field()
	for y := 0; y < SCREEN_HEIGHT; y++ {
		for x := 0; x < SCREEN_WIDTH; x++ {
			c := f.Cell(x, y)
			if c.Alive {
				l := 0.9
				if c.Count == 2 {
					l = 1.0
				}
				c := colorful.Hcl(float64(c.Hue), 1.3, l)
				s.Set(x, y, color.RGBA{uint8(c.R * 255), uint8(c.G * 255), uint8(c.B * 255), 128})
			} else {
				s.Set(x, y, color.RGBA{0, 0, 0, 128})
			}
		}
	}
}

func ScreenToImage(s *Screen) image.Image {
	dst := image.NewRGBA(s.Bounds())
	// draw.DrawMask(dst, dst.Bounds(), s, image.ZP, s, image.ZP, draw.Over)
	draw.Draw(dst, dst.Bounds(), s, image.ZP, draw.Over)
	return dst
}

func ScreenToPalettedImage(s *Screen) *image.Paletted {
	dst := image.NewPaletted(s.Bounds(), palette.Plan9)
	// draw.DrawMask(dst, dst.Bounds(), s, image.ZP, s, image.ZP, draw.Over)
	draw.Draw(dst, dst.Bounds(), s, image.ZP, draw.Over)
	return dst
}

func BlendImages(a, b image.Image, steps int) []*image.Paletted {
	imgs := make([]*image.Paletted, steps)
	for i := 0; i < steps; i++ {
		dst := image.NewPaletted(a.Bounds(), palette.Plan9)
		draw.Draw(dst, dst.Bounds(), a, image.ZP, draw.Over)
		alpha := uint8(255.0 / float32(steps) * float32(i+1))
		draw.DrawMask(dst, dst.Bounds(), b, image.ZP, image.NewUniform(color.Alpha{alpha}), image.ZP, draw.Over)
		imgs[i] = dst
	}
	return imgs
}
