package srv

import (
	"bytes"
	"image"
	"image/draw"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net"
	"time"

	"github.com/nfnt/resize"

	"github.com/ktt-ol/go-insta"
)

func Server(ic insta.Client) {
	l, err := net.Listen("tcp", ":2323")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn) {
			defer c.Close()
			log.Println("got conn")

			// var maxBytes = 1024 * 1024 * 5
			// r := io.LimitReader(c, int64(maxBytes))
			var buf bytes.Buffer
			tee := io.TeeReader(c, &buf)

			config, format, err := image.DecodeConfig(tee)
			log.Println("decode")
			if err != nil {
				c.Write([]byte(err.Error()))
				log.Print(err)
				return
			}

			if config.Height*config.Width > 10000000 {
				c.Write([]byte("too large"))
				log.Print(err)
				return
			}

			r := io.MultiReader(&buf, c)

			if format == "gif" {
				g, err := gif.DecodeAll(r)
				log.Println("decode all")
				if err != nil {
					c.Write([]byte(err.Error()))
					log.Print(err)
					return
				}
				scr := insta.NewScreen()
				bnds := image.Rect(0, 0, insta.ScreenWidth, insta.ScreenHeight)
				for i := range g.Image {
					img := resize.Resize(insta.ScreenWidth, insta.ScreenHeight, g.Image[i], resize.Bilinear)
					draw.Draw(scr, bnds, img, image.ZP, draw.Over)
					ic.SetScreenImmediate(scr)
					time.Sleep(time.Duration(g.Delay[i]) * time.Millisecond)
				}
			} else {
				img, _, err := image.Decode(r)
				if err != nil {
					c.Write([]byte(err.Error()))
					log.Print(err)
					return
				}

				img = resize.Resize(insta.ScreenWidth, insta.ScreenHeight, img, resize.Bilinear)

				scr := insta.NewScreen()
				bnds := image.Rect(0, 0, insta.ScreenWidth, insta.ScreenHeight)
				draw.Draw(scr, bnds, img, image.ZP, draw.Over)

				ic.SetScreenImmediate(scr)
			}
			// Shut down the connection.
		}(conn)
	}
}
