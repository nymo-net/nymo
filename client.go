package nymo

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/http3"
	"github.com/nymo-net/nymo/pb"
	"golang.org/x/net/http2"
	"google.golang.org/protobuf/proto"
)

func (u *user) DialPeer(addr string) (*peer, error) {
	tlsConfig := &tls.Config{InsecureSkipVerify: true}

	var r http.RoundTripper
	var material []byte
	var setHandshake func()

	switch {
	case strings.HasPrefix(addr, "udp://"):
		r = &http3.RoundTripper{
			TLSClientConfig: tlsConfig,
			Dial: func(_, addr string, tlsCfg *tls.Config, cfg *quic.Config) (quic.EarlySession, error) {
				session, err := quic.DialAddrEarly(addr, tlsCfg, cfg)
				if err != nil {
					return nil, err
				}
				state := session.ConnectionState()
				material, err = state.TLS.ExportKeyingMaterial(nymoName, nil, blockSize)
				if err != nil {
					return nil, err
				}
				setHandshake()
				return session, err
			},
		}
	case strings.HasPrefix(addr, "tcp://"):
		r = &http2.Transport{
			TLSClientConfig: tlsConfig,
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				client, err := tls.Dial(network, addr, cfg)
				if err != nil {
					return nil, err
				}
				state := client.ConnectionState()
				material, err = state.ExportKeyingMaterial(nymoName, nil, blockSize)
				if err != nil {
					return nil, err
				}
				setHandshake()
				return client, nil
			},
		}
	default:
		return nil, fmt.Errorf("%s: unknown address format", addr)
	}

	request, err := http.NewRequest(http.MethodPost, "https"+addr[3:], nil)
	if err != nil {
		return nil, err
	}

	reader, writer := io.Pipe()
	request.Body = reader

	setHandshake = func() {
		handshake := pb.PeerHandshake{
			Version:    nymoVersion,
			Cohort:     u.cohort,
			Pow:        calcPoW(material),
			PeerTokens: nil, // TODO
		}

		marshal, err := proto.Marshal(&handshake)
		if err != nil {
			panic(err)
		}

		request.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString(marshal))
	}

	resp, err := r.RoundTrip(request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, resp.Body.Close()
	}

	return u.NewPeerAsClient(request.Context(), resp.Body, writer, material)
}
