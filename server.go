package nymo

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/lucas-clemente/quic-go/http3"
	"github.com/nymo-net/nymo/pb"
	"google.golang.org/protobuf/proto"
)

type server struct {
	user *user
}

func (s *server) validate(r *http.Request) (*pb.PeerHandshake, []byte) {
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
	handshake, sessionKey := s.validate(r)
	if handshake == nil {
		http.DefaultServeMux.ServeHTTP(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	p, err := s.user.NewPeerAsServer(r.Context(), r.Body, &writeFlusher{w, w.(http.Flusher)}, handshake, sessionKey)
	if err != nil {
		s.user.cfg.Logger.Println(err)
		http.DefaultServeMux.ServeHTTP(w, r)
		return
	}

	<-p.ctx.Done()
}

func (u *user) RunServer(listenAddr, certFile, keyFile string) error {
	var serveFunc func(addr, certFile, keyFile string, handler http.Handler) error
	switch {
	case strings.HasPrefix(listenAddr, "udp://"):
		serveFunc = http3.ListenAndServeQUIC
	case strings.HasPrefix(listenAddr, "tcp://"):
		serveFunc = http.ListenAndServeTLS
	default:
		return fmt.Errorf("%s: unknown address format", listenAddr)
	}
	return serveFunc(listenAddr[6:], certFile, keyFile, &server{user: u})
}
