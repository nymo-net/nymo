package nymo

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/lucas-clemente/quic-go/http3"
	"github.com/nymo-net/nymo/pb"
	"google.golang.org/protobuf/proto"
)

type server struct {
	lock  sync.RWMutex
	peers []*Peer
}

func (s *server) validate(r *http.Request) *pb.PeerHandshake {
	if r.Method != http.MethodPost {
		return nil
	}
	if r.URL.Path != "/" {
		return nil
	}

	auth := r.Header.Get("Authorization")
	const prefix = "Basic "
	if !strings.HasPrefix(auth, prefix) {
		return nil
	}

	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return nil
	}

	peer := new(pb.PeerHandshake)
	err = proto.Unmarshal(c, peer)
	if err != nil {
		return nil
	}

	if validatePoW([]byte("localhost:443"), peer.Pow) {
		return peer
	}
	return nil
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handshake := s.validate(r)
	if handshake == nil {
		return
	}

	w.WriteHeader(http.StatusSwitchingProtocols)
	peer, err := NewPeerAsServer(w.(http3.DataStreamer).DataStream(), handshake)
	if err != nil {
		log.Println(err)
		return
	}

	s.lock.Lock()
	s.peers = append(s.peers, peer)
	s.lock.Unlock()
}

func RunServer(listenAddr, certFile, keyFile string) error {
	return http3.ListenAndServeQUIC(listenAddr, certFile, keyFile, new(server))
}
