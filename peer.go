package nymo

import (
	"io"

	"github.com/nymo-net/nymo/pb"
)

type Peer struct {
	conn   io.ReadWriteCloser
	cohort uint32
	tokens [][]byte
}

func NewPeerAsServer(conn io.ReadWriteCloser, handshake *pb.PeerHandshake) (*Peer, error) {
	if err := sendMessage(conn, &pb.HandshakeOK{
		Cohort:     0,
		PeerTokens: nil, // TODO: add peer tokens
	}); err != nil {
		return nil, err
	}

	return &Peer{
		conn:   conn,
		cohort: handshake.Cohort,
		tokens: handshake.PeerTokens,
	}, nil
}

func NewPeerAsClient(conn io.ReadWriteCloser) (*Peer, error) {
	var ok pb.HandshakeOK
	if err := recvMessage(conn, &ok); err != nil {
		return nil, err
	}

	return &Peer{
		conn:   conn,
		cohort: ok.Cohort,
		tokens: ok.PeerTokens,
	}, nil
}
