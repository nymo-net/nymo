package nymo

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
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

var noMTlsAsked = errors.New("server did not properly ask for mTLS")
var peerConnected = errors.New("peer connected")

func (u *User) dialNewPeers(ctx context.Context) {
	u.peerCleanup()

	func() {
		enum := u.db.EnumeratePeers()
		defer enum.Close()

		var outErr error
		for u.shouldConnectPeers() && enum.Next(outErr) {
			url := enum.Url()
			if u.retry.noRetry(url) {
				outErr = nil
				continue
			}

			outErr = u.dialPeer(ctx, enum)
			if ctx.Err() != nil {
				return
			}
			if outErr != nil {
				u.retry.add(url, u.cfg.PeerRetryTime)
			}
			if errors.Is(outErr, peerConnected) {
				outErr = nil
			}
		}
	}()

	u.peerLock.RLock()
	defer u.peerLock.RUnlock()

	// ask for more in-cohort peers
	if u.numIn < u.cfg.MaxInCohortConn {
		in := u.numIn
		for _, p := range u.peers {
			if ctx.Err() != nil {
				return
			}
			if p != nil && p.requestPeer(u.peerSameCohort) {
				in++
				if in >= u.cfg.MaxInCohortConn {
					break
				}
			}
		}
	}

	// ask for more out-of-cohort peers
	out := u.total - u.numIn
	if out < u.cfg.MaxOutCohortConn {
		for _, p := range u.peers {
			if ctx.Err() != nil {
				return
			}
			if p != nil && p.requestPeer(func(c uint32) bool { return !u.peerSameCohort(c) }) {
				out++
				if out >= u.cfg.MaxOutCohortConn {
					break
				}
			}
		}
	}
}

func (u *User) dialPeer(ctx context.Context, handle PeerEnumerate) error {
	reserver := u.reserveServer(handle.Cohort())
	if reserver == nil {
		return nil
	}
	defer reserver.rollback()

	var askedForHandshake bool
	var r http.RoundTripper
	var material []byte
	var peerId [hashTruncate]byte
	var setHandshake func()

	tlsConfig := &tls.Config{
		InsecureSkipVerify: !u.cfg.VerifyServerCert,
		GetClientCertificate: func(info *tls.CertificateRequestInfo) (*tls.Certificate, error) {
			if len(info.AcceptableCAs) > 0 {
				return nil, noMTlsAsked
			}
			askedForHandshake = true
			return &u.cert, nil
		},
	}

	validateState := func(state tls.ConnectionState) (err error) {
		if !askedForHandshake {
			return noMTlsAsked
		}
		id := hasher(state.PeerCertificates[0].Raw)
		peerId = truncateHash(id[:])
		if !reserver.reserveId(&peerId) {
			return peerConnected
		}
		material, err = state.ExportKeyingMaterial(nymoName, nil, hashSize)
		if err != nil {
			return err
		}
		setHandshake()
		return nil
	}

	addr := handle.Url()
	switch {
	case strings.HasPrefix(addr, "udp://"):
		r = &http3.RoundTripper{
			TLSClientConfig: tlsConfig,
			Dial: func(_, addr string, tlsCfg *tls.Config, cfg *quic.Config) (quic.EarlySession, error) {
				session, err := quic.DialAddrEarly(addr, tlsCfg, cfg)
				if err != nil {
					return nil, err
				}
				err = validateState(session.ConnectionState().TLS.ConnectionState)
				if err != nil {
					_ = session.CloseWithError(0, "")
					return nil, err
				}
				return session, nil
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
				err = validateState(client.ConnectionState())
				if err != nil {
					_ = client.Close()
					return nil, err
				}
				return client, nil
			},
		}
	default:
		return fmt.Errorf("%s: unknown address format", addr)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, "https"+addr[3:], nil)
	if err != nil {
		return err
	}

	reader, writer := io.Pipe()
	request.Body = reader

	setHandshake = func() {
		handshake := pb.PeerHandshake{
			Version: nymoVersion,
			Pow:     calcPoW(material),
		}

		marshal, err := proto.Marshal(&handshake)
		if err != nil {
			panic(err)
		}

		request.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString(marshal))
	}

	resp, err := r.RoundTrip(request)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return resp.Body.Close()
	}

	p, err := u.newPeerAsClient(request.Context(), handle, resp.Body, writer, peerId, material)
	if err != nil {
		return err
	}
	reserver.commit(p)
	return nil
}
