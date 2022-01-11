package nymo

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"strings"

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
		http.DefaultServeMux.ServeHTTP(w, r)
		return
	}

	reserver.commit(p)
	<-p.ctx.Done()
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
