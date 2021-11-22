package nymo

import (
	"io"
	"time"

	"github.com/nymo-net/nymo/pb"
)

type Peer struct {
	reader io.ReadCloser
	writer io.Writer
	cohort uint32
	tokens [][]byte
}

func (p *Peer) Ping() {
	for {
		_, err := io.WriteString(p.writer, "pingpingping\n")
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second * 1)
	}
}

func NewPeerAsServer(r io.ReadCloser, w io.Writer, handshake *pb.PeerHandshake) (*Peer, error) {
	if err := sendMessage(w, &pb.HandshakeOK{
		Cohort:     0,
		PeerTokens: nil, // TODO: add peer tokens
	}); err != nil {
		return nil, err
	}

	return &Peer{
		reader: r,
		writer: w,
		cohort: handshake.Cohort,
		tokens: handshake.PeerTokens,
	}, nil
}

func NewPeerAsClient(r io.ReadCloser, w io.Writer) (*Peer, error) {
	var ok pb.HandshakeOK
	if err := recvMessage(r, &ok); err != nil {
		return nil, err
	}

	return &Peer{
		reader: r,
		writer: w,
		cohort: ok.Cohort,
		tokens: ok.PeerTokens,
	}, nil
}
