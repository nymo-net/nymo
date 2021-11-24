package nymo

import (
	"io"

	"github.com/nymo-net/nymo/pb"
)

type peer struct {
	reader io.ReadCloser
	writer io.Writer
	cohort uint32
	tokens [][]byte
	key    []byte
}

func (u *user) NewPeerAsServer(r io.ReadCloser, w io.Writer, handshake *pb.PeerHandshake, sKey []byte) (*peer, error) {
	if err := sendMessage(w, &pb.HandshakeOK{
		Cohort:     u.cohort,
		PeerTokens: nil, // TODO: add peer tokens
	}); err != nil {
		return nil, err
	}

	return &peer{
		reader: r,
		writer: w,
		cohort: handshake.Cohort,
		tokens: handshake.PeerTokens,
		key:    sKey,
	}, nil
}

func (u *user) NewPeerAsClient(r io.ReadCloser, w io.Writer, sKey []byte) (*peer, error) {
	var ok pb.HandshakeOK
	if err := recvMessage(r, &ok); err != nil {
		return nil, err
	}

	return &peer{
		reader: r,
		writer: w,
		cohort: ok.Cohort,
		tokens: ok.PeerTokens,
		key:    sKey,
	}, nil
}
