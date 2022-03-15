package nymo

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/lucas-clemente/quic-go/http3"
	"github.com/nymo-net/nymo/pb"
	"google.golang.org/protobuf/proto"
)

type server struct {
	user *User
}

func validate(r *http.Request) (*pb.PeerHandshake, []byte) {
	if r.Method != http.MethodPost {
		return nil, nil
	}
	if r.URL.Path != "/" {
		return nil, nil
	}

	auth := r.Header.Get("Authorization")
	const prefix = "Basic "
	if !strings.HasPrefix(auth, prefix) {
		return nil, nil
	}

	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return nil, nil
	}

	p := new(pb.PeerHandshake)
	err = proto.Unmarshal(c, p)
	if err != nil {
		return nil, nil
	}

	if p.Version != nymoVersion {
		return nil, nil
	}

	material, err := r.TLS.ExportKeyingMaterial(nymoName, nil, blockSize)
	if err != nil {
		return nil, nil
	}

	if validatePoW(material, p.Pow) {
		return p, material
	}
	return nil, nil
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(r.TLS.PeerCertificates) <= 0 {
		http.DefaultServeMux.ServeHTTP(w, r)
		return
	}

	certHash := hasher(r.TLS.PeerCertificates[0].Raw)
	reserver := s.user.reserveClient(truncateHash(certHash[:]))
	if reserver == nil {
		w.WriteHeader(http.StatusTeapot)
		return
	}
	defer reserver.rollback()

	handshake, material := validate(r)
	if handshake == nil {
		http.DefaultServeMux.ServeHTTP(w, r)
		return
	}
	if !reserver.reserveCohort(handshake.Cohort) {
		w.WriteHeader(http.StatusTeapot)
		return
	}
	handle := s.user.db.ClientHandle(certHash[:hashTruncate])

	w.WriteHeader(http.StatusOK)
	p, err := s.user.newPeerAsServer(
		r.Context(), handle, r.Body, &writeFlusher{w, w.(http.Flusher)}, material, handshake.Cohort)
	if err != nil {
		handle.Disconnect(err)
		return
	}

	reserver.commit(p)
	<-p.ctx.Done()
}

func (u *User) RunServerUpnp(ctx context.Context, listenAddr string) error {
	var protocol string
	switch {
	case strings.HasPrefix(listenAddr, "udp://"):
		protocol = "UDP"
	case strings.HasPrefix(listenAddr, "tcp://"):
		protocol = "TCP"
	default:
		return fmt.Errorf("%s: unsupported address format", listenAddr)
	}

	addr := listenAddr[6:]
	host, portS, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}
	port, err := strconv.Atoi(portS)
	if err != nil {
		return err
	}
	if port <= 0 || port > math.MaxUint16 {
		return errors.New("port number out of range")
	}
	portS = strconv.Itoa(port)
	portU := uint16(port)

	client, err := pickRouterClient(ctx)
	if err != nil {
		return err
	}
	extAddr, err := client.GetExternalIPAddressCtx(ctx)
	if err != nil {
		return err
	}

	err = client.AddPortMappingCtx(ctx, "", portU, protocol, portU, host, true, nymoName, 3600)
	if err != nil {
		return err
	}
	defer client.DeletePortMappingCtx(context.Background(), "", portU, protocol)

	ctx, cancel := context.WithCancel(ctx)
	srv := &http.Server{
		Handler: &server{user: u},
		Addr:    addr,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{u.cert},
			ClientAuth:   tls.RequestClientCert,
		},
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
		ErrorLog: u.cfg.Logger,
	}

	go func() {
		func() {
			defer cancel()

			ticker := time.NewTicker(time.Second * 3000)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					err = client.AddPortMappingCtx(ctx, "", portU, protocol, portU, host, true, nymoName, 3600)
					if err != nil {
						u.cfg.Logger.Print(err)
						return
					}
				}
			}
		}()
		_ = srv.Shutdown(context.Background())
	}()

	switch protocol {
	case "UDP":
		listenAddr = "udp://" + net.JoinHostPort(extAddr, portS)
	case "TCP":
		listenAddr = "tcp://" + net.JoinHostPort(extAddr, portS)
	}

	hash := hasher([]byte(listenAddr))

	switch protocol {
	case "UDP":
		u.retry.addSelf(listenAddr)
		u.db.AddPeer(listenAddr, &pb.Digest{
			Hash:   hash[:hashTruncate],
			Cohort: u.cohort,
		})

		return (&http3.Server{Server: srv}).ListenAndServe()
	case "TCP":
		u.retry.addSelf(listenAddr)
		u.db.AddPeer(listenAddr, &pb.Digest{
			Hash:   hash[:hashTruncate],
			Cohort: u.cohort,
		})

		return srv.ListenAndServeTLS("", "")
	}
	panic("not reachable")
}

func (u *User) RunServer(ctx context.Context, listenAddr string) error {
	srv := &http.Server{
		Handler: &server{user: u},
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{u.cert},
			ClientAuth:   tls.RequestClientCert,
		},
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
		ErrorLog: u.cfg.Logger,
	}
	go func() {
		<-ctx.Done()
		_ = srv.Shutdown(context.Background())
	}()

	hash := hasher([]byte(listenAddr))

	switch {
	case strings.HasPrefix(listenAddr, "udp://"):
		u.retry.addSelf(listenAddr)
		u.db.AddPeer(listenAddr, &pb.Digest{
			Hash:   hash[:hashTruncate],
			Cohort: u.cohort,
		})

		srv.Addr = listenAddr[6:]
		return (&http3.Server{Server: srv}).ListenAndServe()
	case strings.HasPrefix(listenAddr, "tcp://"):
		u.retry.addSelf(listenAddr)
		u.db.AddPeer(listenAddr, &pb.Digest{
			Hash:   hash[:hashTruncate],
			Cohort: u.cohort,
		})

		srv.Addr = listenAddr[6:]
		return srv.ListenAndServeTLS("", "")
	default:
		return fmt.Errorf("%s: unknown address format", listenAddr)
	}
}

func (u *User) ListServers() (ret []string) {
	u.retry.l.Lock()
	defer u.retry.l.Unlock()

	for k, v := range u.retry.m {
		if v == emptyTime {
			ret = append(ret, k)
		}
	}

	return
}
