package insta

import (
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"math"
	"strconv"
)

const (
	PanelWidth  = 18
	PanelHWidth = PanelWidth / 2
	PanelHeight = 18

	PanelsX = 3
	PanelsY = 2

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

func (s *Screen) Copy() *Screen {
	r := NewScreen()
	copy(r.Pix, s.Pix)
	return r
}

func (s *Screen) Set(x, y int, c color.Color) {
	if x < 0 || y < 0 || x >= ScreenWidth || y >= ScreenHeight {
		return
	}
	rgb := color.RGBA{}
	if rgba, ok := c.(color.RGBA); ok {
		rgb = rgba
	} else {
		r, g, b, _ := c.RGBA()
		rgb.R = uint8(r / 256)
		rgb.G = uint8(g / 256)
		rgb.B = uint8(b / 256)
	}
	offset := y*LineStride + x*PixelStride
	s.Pix[offset] = rgb.R
	s.Pix[offset+1] = rgb.G
	s.Pix[offset+2] = rgb.B
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
		copy(l[py*PanelHWidth*PixelStride:], s.Pix[offset:offset+PanelHWidth*PixelStride])
		copy(r[py*PanelHWidth*PixelStride:], s.Pix[offset+PanelHWidth*PixelStride:offset+PanelWidth*PixelStride])
	}
	return l, r
}

func HsvToRgb(h, s, v float64) (r, g, b float64) {
	var i int
	var f, p, q, t float64

	if s == 0 {
		// achromatic (grey)
		r = v
		g = v
		b = v
		return
	}

	h = h / 60 // sector 0 to 5
	i = int(math.Floor(h))
	f = h - float64(i) // factorial part of h
	p = v * (1 - s)
	q = v * (1 - s*f)
	t = v * (1 - s*(1-f))

	switch i {
	case 0:
		r = v
		g = t
		b = p
	case 1:
		r = q
		g = v
		b = p
	case 2:
		r = p
		g = v
		b = t
	case 3:
		r = p
		g = q
		b = v
	case 4:
		r = t
		g = p
		b = v
	default: // case 5:
		r = v
		g = p
		b = q
	}
	return
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

func BlendScreens(a, b *Screen, steps int) []*Screen {
	screens := make([]*Screen, steps)
	if steps <= 1 {
		return []*Screen{b}
	}
	for i := 0; i < steps; i++ {
		dst := NewScreen()

		t := 1 / float32(steps) * float32(i+1)
		for i := 0; i < ScreenWidth*ScreenHeight*PixelStride; i++ {
			dst.Pix[i] = uint8((1-t)*float32(a.Pix[i]) + t*float32(b.Pix[i]))
		}
		screens[i] = dst
	}
	return screens
}
