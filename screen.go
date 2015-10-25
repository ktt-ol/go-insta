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
	PanelWidth  = 18
	PanelHwidth = PanelWidth / 2
	PanelHeight = 18

	PanelsX = 4
	PanelsY = 3

	ScreenWidth  = PanelWidth * PanelsX
	ScreenHeight = PanelHeight * PanelsY

	PixelStride  = 3
	PanelStride  = PixelStride * PanelWidth
	LineStride   = PanelStride * PanelsX
	PanelYStride = LineStride * PanelHeight
)

type Screen struct {
	Pix []uint8
}

func NewScreen() *Screen {
	return &Screen{
		Pix: make([]uint8, ScreenWidth*ScreenHeight*PixelStride),
	}
}

func (s *Screen) Set(x, y int, c color.RGBA) {
	offset := y*LineStride + x*PixelStride
	s.Pix[offset] = c.R
	s.Pix[offset+1] = c.G
	s.Pix[offset+2] = c.B
}

func (s *Screen) At(x, y int) color.Color {
	offset := y*LineStride + x*PixelStride
	return color.RGBA{
		s.Pix[offset],
		s.Pix[offset+1],
		s.Pix[offset+2],
		255,
	}
}

func (s *Screen) Bounds() image.Rectangle {
	return image.Rect(0, 0, ScreenWidth, ScreenHeight)
}

func (s *Screen) ColorModel() color.Model {
	return color.RGBAModel
}

func (s *Screen) String() string {
	var b []byte
	offset := 0
	for y := 0; y < ScreenHeight; y++ {
		for x := 0; x < ScreenWidth; x++ {
			for c := 0; c < PixelStride; c++ {
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

	panelOffset := x*PanelStride + y*PanelYStride
	for py := 0; py < PanelHeight; py++ {
		offset := panelOffset + py*LineStride
		copy(l[py*PanelWidth/2:], s.Pix[offset:offset+PanelHwidth*PixelStride])
		copy(r[py*PanelWidth/2:], s.Pix[offset+PanelHwidth:offset+PanelWidth*PixelStride])

	}
	return l, r
}

func LifeToScreen(l *life.Life, s *Screen) {
	f := l.Field()
	for y := 0; y < ScreenHeight; y++ {
		for x := 0; x < ScreenWidth; x++ {
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
