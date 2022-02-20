package nymo

import (
	"bufio"
	"bytes"
	"context"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/nymo-net/nymo/pb"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

const (
	multicastGroupIpv4 = "239.192.152.143:6770"
	multicastGroupIpv6 = "[ff15::efc0:988f]:6770"

	orgLocalTTL  = 31
	announceTime = 5 * time.Minute

	maxUdpPayloadSize = 576 - 60 - 8
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (u *User) peerAnnounce(ctx context.Context, c net.Conn, host string) error {
	ticker := time.NewTicker(announceTime)
	defer ticker.Stop()

	for {
		srvs := u.ListServers()
		if len(srvs) > 0 {
			srv := srvs[rand.Intn(len(srvs))]
			req := http.Request{
				Method: "BT-SEARCH",
				Host:   host,
				URL:    &url.URL{Opaque: "*"},
				Header: http.Header{
					"User-Agent": nil,
					"Port":       []string{strconv.Itoa(int(u.cohort))},
					"Infohash":   []string{srv},
				},
			}
			err := req.Write(c)
			if err != nil {
				return err
			}
		}
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}
	}
}

func (u *User) ipv4PeerAnnounce(ctx context.Context) error {
	var dialer net.Dialer

	c, err := dialer.DialContext(ctx, "udp4", multicastGroupIpv4)
	if err != nil {
		return err
	}
	defer c.Close()

	err = ipv4.NewPacketConn(c.(net.PacketConn)).SetMulticastTTL(orgLocalTTL)
	if err != nil {
		return err
	}

	return u.peerAnnounce(ctx, c, multicastGroupIpv4)
}

func (u *User) ipv6PeerAnnounce(ctx context.Context) error {
	var dialer net.Dialer

	c, err := dialer.DialContext(ctx, "udp6", multicastGroupIpv6)
	if err != nil {
		return err
	}
	defer c.Close()

	err = ipv6.NewPacketConn(c.(net.PacketConn)).SetMulticastHopLimit(orgLocalTTL)
	if err != nil {
		return err
	}

	return u.peerAnnounce(ctx, c, multicastGroupIpv6)
}

func parsePacket(p []byte, host string) ([]string, uint32) {
	req, err := http.ReadRequest(bufio.NewReader(bytes.NewBuffer(p)))
	if err != nil || req.Host != host ||
		req.Method != "BT-SEARCH" || req.RequestURI != "*" || req.Proto != "HTTP/1.1" {
		return nil, 0
	}
	addrs := req.Header["Infohash"]
	cohort := req.Header["Port"]
	if len(cohort) != 1 || len(addrs) <= 0 {
		return nil, 0
	}
	c, err := strconv.Atoi(cohort[0])
	if err != nil {
		return nil, 0
	}
	if c < 0 || c >= cohortNumber {
		return nil, 0
	}
	return addrs, uint32(c)
}

func (u *User) peerDiscover(ctx context.Context, host string) error {
	addr, err := net.ResolveUDPAddr("udp", host)
	if err != nil {
		return err
	}

	l, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	go func() {
		defer l.Close()
		<-ctx.Done()
	}()

	buf := make([]byte, maxUdpPayloadSize)
	for {
		pu, err := l.Read(buf)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return err
		}
		addrs, cohort := parsePacket(buf[:pu], host)
		for _, a := range addrs {
			hash := hasher([]byte(a))
			u.db.AddPeer(a, &pb.Digest{
				Hash:   hash[:hashTruncate],
				Cohort: cohort,
			})
		}
	}
}

func (u *User) ipv4PeerDiscover(ctx context.Context) error {
	return u.peerDiscover(ctx, multicastGroupIpv4)
}

func (u *User) ipv6PeerDiscover(ctx context.Context) error {
	return u.peerDiscover(ctx, multicastGroupIpv6)
}
