package nymo

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
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
	defer reserver.cleanup(s.user, p)
	<-p.ctx.Done()
}

func (u *User) RunServerUpnp(ctx context.Context, serverAddr string) error {
	var protocol string
	switch {
	case strings.HasPrefix(serverAddr, "udp://"):
		protocol = "UDP"
	case strings.HasPrefix(serverAddr, "tcp://"):
		protocol = "TCP"
	default:
		return fmt.Errorf("%s: unsupported address format", serverAddr)
	}

	addr := serverAddr[6:]
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

	go func() {
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

	switch protocol {
	case "UDP":
		serverAddr = "udp://" + net.JoinHostPort(extAddr, portS)
	case "TCP":
		serverAddr = "tcp://" + net.JoinHostPort(extAddr, portS)
	}

	return u.RunServer(ctx, serverAddr, addr)
}

func (u *User) RunServer(ctx context.Context, serverAddr, listenAddr string) error {
	srv := &http.Server{
		Addr:    listenAddr,
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
	var cl io.Closer = srv
	go func() {
		<-ctx.Done()
		_ = cl.Close()
	}()

	hash := hasher([]byte(serverAddr))
	switch {
	case strings.HasPrefix(serverAddr, "udp://"):
		u.retry.addSelf(serverAddr)
		u.db.AddPeer(serverAddr, &pb.Digest{
			Hash:   hash[:hashTruncate],
			Cohort: u.cohort,
		})
		s := &http3.Server{Server: srv}
		cl = s
		return s.ListenAndServe()
	case strings.HasPrefix(serverAddr, "tcp://"):
		u.retry.addSelf(serverAddr)
		u.db.AddPeer(serverAddr, &pb.Digest{
			Hash:   hash[:hashTruncate],
			Cohort: u.cohort,
		})
		return srv.ListenAndServeTLS("", "")
	default:
		return fmt.Errorf("%s: unknown address format", serverAddr)
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
