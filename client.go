package insta

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/draw"
	"log"
	"net"
	"time"

	"fmt"
)

const (
	syncPort = 8500
	dataPort = 9410
)

type pkt struct {
	head            [34]byte
	panelLeft       [486]uint8
	brightnessLeft  uint8
	contrastLeft    uint8
	afterglowLeft   uint8
	unknownLeft     uint8
	trailLeft       [49]byte
	panelRight      [486]uint8
	brightnessRight uint8
	contrastRight   uint8
	afterglowRight  uint8
	unknownRight    uint8
	trailRight      [49]byte
}

func newPkt() *pkt {
	p := pkt{}
	copy(p.head[:], []byte("INSTA-INET\x00\x01\x00\x00\x01\xac\x10\x05\x00\x00"+
		/* img 1, sync 8: */ "\x01"+"\x00\x01"+
		/* frame counter */ "\x00\x00"+"\x00\x00\x00\x00\x01\xe6\x00\x1a\x00"))
	copy(p.trailLeft[:], []byte("\xff\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x14\x91\x49\x45\x01"+
		/* 00|01|02 cycle */ "\x00"+"\x00\x01"+
		/* e3|6c|e0 cycle */ "\xe3"+"\x80\x16\x08"+
		/* "checksum" (cycle+ip) */ "\x00\x00\x00\x00\x00\x00\x00\x00\x00"+"\x00\x00\x00\x00\x00\x00"))
	copy(p.trailRight[:], p.trailLeft[:])
	p.brightnessLeft = 128
	p.brightnessRight = 128
	p.contrastLeft = 128
	p.contrastRight = 128
	p.afterglowLeft = 80
	p.afterglowRight = 80
	p.unknownLeft = 0
	p.unknownRight = 0
	return &p
}

type syncPkt struct {
	head [27]byte
}

func newSyncPkt() *syncPkt {
	s := syncPkt{}
	copy(s.head[:], []byte("INSTA-INET\x00\x01\x00\x00\x01\xac\x10\x05\x00\x00"+
		/* img 1, sync 8: */ "\x08"+"\x00\x00\x00\x00\x00\x00"))
	return &s
}

type Client interface {
	SetScreen(s *Screen)
	SetFPS(int)
	Run()
	SetAfterglow(float64)
}

type InstaClient struct {
	syncSock   *net.UDPConn
	dataSock   *net.UDPConn
	panelAddrs []*net.UDPAddr
	syncBytes  []byte
	dataPkt    *pkt
	imgs       chan image.Image
	fps        int
}

func NewInstaClient(addrs []string) (*InstaClient, error) {
	c := InstaClient{imgs: make(chan image.Image, 1), fps: 50}
	if len(addrs) != PanelsX*PanelsY {
		return nil, fmt.Errorf("invalid number of addresses, got %d for %dx%d panels",
			len(addrs), PanelsX, PanelsY)
	}
	for _, addr := range addrs {
		udpAddr := &net.UDPAddr{IP: net.ParseIP(addr), Port: dataPort}
		c.panelAddrs = append(c.panelAddrs, udpAddr)
	}
	var err error
	c.dataSock, err = net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	if err != nil {
		return nil, err
	}

	lip, err := localSourceIP(net.ParseIP(addrs[0]))
	if err != nil {
		return nil, err
	}
	c.syncSock, err = net.DialUDP("udp4",
		&net.UDPAddr{IP: lip, Port: syncPort},
		&net.UDPAddr{IP: net.IPv4(255, 255, 255, 255), Port: dataPort})
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	syncPkt := newSyncPkt()
	if err := binary.Write(buf, binary.LittleEndian, syncPkt); err != nil {
		return nil, err
	}
	c.syncBytes = buf.Bytes()
	c.dataPkt = newPkt()
	return &c, nil
}

func (c *InstaClient) SetScreen(s *Screen) {
	select {
	case c.imgs <- s.Copy():
	default: // skip screen
	}
}

func (c *InstaClient) SetAfterglow(v float64) {
	if v < 0 {
		v = 0
	}
	if v > 1 {
		v = 1
	}
	c.dataPkt.afterglowLeft = uint8(v * 255)
	c.dataPkt.afterglowRight = uint8(v * 255)
}

func (c *InstaClient) send(img image.Image) error {
	scr := NewScreen() // TODO, create image2panel function
	draw.Draw(scr, scr.Bounds(), img, image.ZP, draw.Over)

	i := 0
	buf := new(bytes.Buffer)
	for y := 0; y < PanelsY; y++ {
		for x := 0; x < PanelsX; x++ {
			c.dataPkt.panelLeft, c.dataPkt.panelRight = scr.Panel(x, y)
			err := binary.Write(buf, binary.LittleEndian, c.dataPkt)
			if err != nil {
				return err
			}
			c.dataSock.WriteTo(buf.Bytes(), c.panelAddrs[i])
			// n, err := c.dataSock.WriteTo(buf.Bytes(), c.panelAddrs[i])
			// if err != nil {
			// 	return err
			// }
			// if n != buf.Len() {
			// 	return fmt.Errorf("not all bytes sent: %d of %d", n, buf.Len())
			// }
			buf.Reset()
			i += 1
		}
	}

	return nil
}

func (c *InstaClient) SetFPS(fps int) {
	if fps < 0 || fps > 100 {
		fps = 50
	}
	c.fps = fps
}

func (c *InstaClient) sync() error {
	n, err := c.syncSock.Write(c.syncBytes)
	if err != nil {
		return err
	}
	if n != len(c.syncBytes) {
		return fmt.Errorf("not all bytes sent: %d of %d", n, len(c.syncBytes))
	}
	return nil
}

func (c *InstaClient) Run() {
	dur := time.Duration(1000/float64(c.fps)) * time.Millisecond
	start := time.Now()
	for img := range c.imgs {
		sendStart := time.Now()
		if err := c.send(img); err != nil {
			log.Printf("error: while sending packages: %v", err)
			time.Sleep(time.Second)
		} else {
			end := time.Now()
			wait := time.Duration(dur.Nanoseconds()-end.Sub(start).Nanoseconds()) * time.Nanosecond
			fmt.Println(start, end, wait, end.Sub(start), sendStart.Sub(start))
			time.Sleep(wait)
			if err := c.sync(); err != nil {
				log.Printf("error: while sending packages: %v", err)
			}
		}
		start = time.Now()
	}
}

// localSourceIP returns the local IP address that is in the same network as target
func localSourceIP(target net.IP) (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return net.IPv4zero, err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return net.IPv4zero, err
		}
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				if v.Contains(target) {
					return v.IP, nil
				}
			}
		}
	}
	return net.IPv4zero, fmt.Errorf("found no interface for %v", target)
}
