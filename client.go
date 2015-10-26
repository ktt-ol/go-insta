package insta

import (
	"bytes"
	"encoding/binary"
	"net"

	"fmt"
)

const (
	SyncPort = 8500
	DataPort = 9410
)

type Pkt struct {
	head            [34]byte
	panelLeft       [486]uint8
	brightnessLeft  uint8
	contrastLeft    uint8
	controlCLeft    uint8
	controlDLeft    uint8
	trailLeft       [49]byte
	panelRight      [486]uint8
	brightnessRight uint8
	contrastRight   uint8
	controlCRight   uint8
	controlDRight   uint8
	trailRight      [49]byte
}

func NewPkt() *Pkt {
	p := Pkt{}
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
	return &p
}

type Sync struct {
	head [27]byte
}

func NewSync() *Sync {
	s := Sync{}
	copy(s.head[:], []byte("INSTA-INET\x00\x01\x00\x00\x01\xac\x10\x05\x00\x00"+
		/* img 1, sync 8: */ "\x08"+"\x00\x00\x00\x00\x00\x00"))
	return &s
}

type Client struct {
	syncSock   *net.UDPConn
	dataSock   *net.UDPConn
	panelAddrs []*net.UDPAddr
	syncPkt    *Sync
	dataPkg    *Pkt
	screen     *Screen
}

func NewClient(addrs []string) (*Client, error) {
	c := Client{}
	if len(addrs) != PanelsX*PanelsY {
		return nil, fmt.Errorf("invalid number of addresses, got %d for %dx%d panels",
			len(addrs), PanelsX, PanelsY)
	}
	for _, addr := range addrs {
		udpAddr := &net.UDPAddr{IP: net.ParseIP(addr), Port: DataPort}
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
		&net.UDPAddr{IP: lip, Port: SyncPort},
		&net.UDPAddr{IP: net.IPv4(255, 255, 255, 255), Port: DataPort})
	if err != nil {
		return nil, err
	}

	c.syncPkt = NewSync()
	c.dataPkg = NewPkt()
	return &c, nil
}

func (c *Client) SetScreen(s *Screen) {
	c.screen = s
}

func (c *Client) Send() error {
	if c.screen == nil {
		return nil
	}

	i := 0
	buf := new(bytes.Buffer)
	for y := 0; y < PanelsY; y++ {
		for x := 0; x < PanelsX; x++ {
			c.dataPkg.panelLeft, c.dataPkg.panelRight = c.screen.Panel(x, y)
			err := binary.Write(buf, binary.LittleEndian, c.dataPkg)
			if err != nil {
				return err
			}
			n, err := c.dataSock.WriteTo(buf.Bytes(), c.panelAddrs[i])
			if err != nil {
				return err
			}
			if n != buf.Len() {
				return fmt.Errorf("not all bytes sent: %d of %d", n, buf.Len())
			}
			buf.Reset()
			i += 1
		}
	}
	if err := binary.Write(buf, binary.LittleEndian, c.syncPkt); err != nil {
		return err
	}

	n, err := c.syncSock.Write(buf.Bytes())
	if err != nil {
		return err
	}
	if n != buf.Len() {
		return fmt.Errorf("not all bytes sent: %d of %d", n, buf.Len())
	}

	return nil
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
