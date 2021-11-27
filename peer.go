package nymo

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/nymo-net/nymo/pb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type peer struct {
	ctx  context.Context
	user *user

	reader io.ReadCloser
	writer io.Writer
	cohort uint32
	tokens [][]byte
	key    []byte

	queue chan proto.Message
}

func (p *peer) Close() {
	p.reader.Close()
}

func (p *peer) recvMessages() {
	defer p.reader.Close()

	for {
		var any anypb.Any
		err := recvMessage(p.reader, &any)
		if err != nil {
			log.Print(err)
			break
		}

		n, err := any.UnmarshalNew()
		if err != nil {
			log.Print(err)
			break
		}

		switch msg := n.(type) {
		case *pb.RequestMsg:
		case *pb.RequestPeer:
		case *pb.MsgList:
		case *pb.MsgContainer:
		case *pb.ResponsePeer:
			_ = msg
		default:
			log.Printf("unknown message type %T", msg)
		}
	}
}

func (p *peer) sendMessages() {
	listTimer := time.NewTicker(p.user.cfg.ListMessageTime)
	defer listTimer.Stop()

	for {
		var msg proto.Message
		select {
		case msg = <-p.queue:
		case <-listTimer.C:
			// TODO
		case <-p.ctx.Done():
			break
		}

		any, err := anypb.New(msg)
		if err != nil {
			log.Panic(err)
		}

		err = sendMessage(p.writer, any)
		if err != nil {
			log.Println(err)
			break
		}
	}
}

func runPeer(p *peer) *peer {
	p.queue = make(chan proto.Message)
	p.user.peers.Store(p, nil)
	go p.sendMessages()
	go p.recvMessages()
	return p
}

func (u *user) NewPeerAsServer(
	ctx context.Context, r io.ReadCloser, w io.Writer, handshake *pb.PeerHandshake, sKey []byte) (*peer, error) {
	if err := sendMessage(w, &pb.HandshakeOK{
		Cohort:     u.cohort,
		PeerTokens: nil, // TODO: add peer tokens
	}); err != nil {
		return nil, err
	}

	return runPeer(&peer{
		ctx:    ctx,
		user:   u,
		reader: r,
		writer: w,
		cohort: handshake.Cohort,
		tokens: handshake.PeerTokens,
		key:    sKey,
	}), nil
}

func (u *user) NewPeerAsClient(ctx context.Context, r io.ReadCloser, w io.Writer, sKey []byte) (*peer, error) {
	var ok pb.HandshakeOK
	if err := recvMessage(r, &ok); err != nil {
		return nil, err
	}

	return runPeer(&peer{
		ctx:    ctx,
		user:   u,
		reader: r,
		writer: w,
		cohort: ok.Cohort,
		tokens: ok.PeerTokens,
		key:    sKey,
	}), nil
}
