package nymo

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/http3"
	"github.com/nymo-net/nymo/pb"
	"golang.org/x/net/http2"
	"google.golang.org/protobuf/proto"
)

func DialPeer(addr string) (*Peer, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	var r http.RoundTripper
	var setHandshake func(*x509.Certificate)

	switch {
	case strings.HasPrefix(addr, "udp://"):
		r = &http3.RoundTripper{
			TLSClientConfig: tlsConfig,
			Dial: func(_, addr string, tlsCfg *tls.Config, cfg *quic.Config) (quic.EarlySession, error) {
				session, err := quic.DialAddrEarly(addr, tlsCfg, cfg)
				if err != nil {
					return nil, err
				}
				setHandshake(session.ConnectionState().TLS.PeerCertificates[0])
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
				setHandshake(client.ConnectionState().PeerCertificates[0])
				return client, nil
			},
		}
	default:
		return nil, fmt.Errorf("%s: unknown address format", addr)
	}

	request, err := http.NewRequest(http.MethodPost, "https"+addr[3:], nil)
	if err != nil {
		log.Panic(err)
	}

	reader, writer := io.Pipe()
	request.Body = reader

	setHandshake = func(cert *x509.Certificate) {
		handshake := pb.PeerHandshake{
			Cohort:     0,
			Pow:        calcPoW(cert.Raw),
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

	return NewPeerAsClient(resp.Body, writer)
}
