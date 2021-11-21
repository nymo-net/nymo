package nymo

import (
	"crypto/tls"
	"encoding/base64"
	"net/http"

	"github.com/lucas-clemente/quic-go/http3"
	"github.com/nymo-net/nymo/pb"
	"google.golang.org/protobuf/proto"
)

func DialPeer(addr string) (*Peer, error) {
	r := &http3.RoundTripper{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	request, err := http.NewRequest(http.MethodPost, "https://"+addr, nil)
	if err != nil {
		panic(err)
	}

	handshake := pb.PeerHandshake{
		Cohort:     0,
		Pow:        calcPoW([]byte(addr)), // TODO
		PeerTokens: nil,
	}
	marshal, err := proto.Marshal(&handshake)
	if err != nil {
		panic(err)
	}
	request.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString(marshal))

	resp, err := r.RoundTrip(request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusSwitchingProtocols {
		return nil, resp.Body.Close()
	}

	return NewPeerAsClient(resp.Body.(http3.DataStreamer).DataStream())
}
