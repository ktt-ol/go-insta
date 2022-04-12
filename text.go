package insta

import (
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/plan9font"
	"golang.org/x/image/math/fixed"
)

func ScrollText(c Client, text string, speed time.Duration, base int) {
	fontWidth := 10
	fontHeight := 20
	fontName := fmt.Sprintf("%dx%d", fontWidth, fontHeight)

	readFile := func(name string) ([]byte, error) {
		return ioutil.ReadFile(filepath.FromSlash(path.Join("plan9fonts/"+fontName, name)))
	}
	fontData, err := readFile(fontName + ".font")
	if err != nil {
		log.Fatal(err)
	}
	face, err := plan9font.ParseFont(fontData, readFile)
	if err != nil {
		log.Fatal(err)
	}

	steps := len(text)*fontWidth + ScreenWidth
	fmt.Println(steps)
	for i := 0; i < steps; i++ {
		scr := NewScreen()
		d := &font.Drawer{
			Dst:  scr,
			Src:  image.White,
			Face: face,
			Dot:  fixed.P(-i+ScreenWidth, base),
		}
		d.DrawString(text)
		c.SetScreen(scr)
		time.Sleep(speed)
	}
	time.Sleep(100 * time.Millisecond)
}
