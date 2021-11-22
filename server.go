package nymo

import (
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/lucas-clemente/quic-go/http3"
	"github.com/nymo-net/nymo/pb"
	"google.golang.org/protobuf/proto"
)

type server struct {
	lock  sync.RWMutex
	peers []*Peer
	cert  []byte
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

	if validatePoW(s.cert, peer.Pow) {
		return peer
	}
	return nil
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handshake := s.validate(r)
	if handshake == nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	peer, err := NewPeerAsServer(r.Body, w, handshake)
	w.(http.Flusher).Flush()
	if err != nil {
		log.Println(err)
		return
	}

	s.lock.Lock()
	s.peers = append(s.peers, peer)
	s.lock.Unlock()

	defer peer.reader.Close()
	io.Copy(os.Stdout, peer.reader)
}

func RunServer(listenAddr, certFile, keyFile string) error {
	certData, err := ioutil.ReadFile(certFile)
	if err != nil {
		return err
	}
	decode, _ := pem.Decode(certData)

	var serveFunc func(addr, certFile, keyFile string, handler http.Handler) error
	switch {
	case strings.HasPrefix(listenAddr, "udp://"):
		serveFunc = http3.ListenAndServeQUIC
	case strings.HasPrefix(listenAddr, "tcp://"):
		serveFunc = http.ListenAndServeTLS
	default:
		return fmt.Errorf("%s: unknown address format", listenAddr)
	}
	return serveFunc(listenAddr[6:], certFile, keyFile, &server{cert: decode.Bytes})
}
