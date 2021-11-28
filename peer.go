package nymo

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/nymo-net/nymo/pb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type peer struct {
	ctx  context.Context
	done uintptr
	url  string

	reader io.ReadCloser
	writer io.Writer
	cohort uint32
	tokens [][]byte
	key    []byte

	queue chan proto.Message
}

func (p *peer) GetToken() []byte {
	if p.url == "" {
		return nil
	}
	hash := hasher([]byte(p.url))
	return hash[8:]
}

func (p *peer) Close() {
	p.reader.Close()
}

func (p *peer) validateTokenPoW(token []byte, pow uint64) bool {
	return validatePoW(append(p.key, token...), pow)
}

func (p *peer) sendProto(msg proto.Message) error {
	select {
	case p.queue <- msg:
		return nil
	case <-p.ctx.Done():
		return p.ctx.Err()
	}
}

func (p *peer) SendPeer(url string) {
	_ = p.sendProto(&pb.ResponsePeer{Address: url})
}

func (p *peer) SendMessage(msg *pb.Message) error {
	// FIXME: demo only, remove later

	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	hash := hasher(data)
	msgCont := &pb.MsgContainer{
		Msg:     data,
		MsgHash: hash[:],
		Pow:     calcPoW(hash[:]),
	}

	return p.sendProto(msgCont)
}

func (u *user) peerConnected(p *peer) {
	token := p.GetToken()
	if token != nil {
		u.db.AddPeer(p.url, token, p.cohort, true)
	}
}

func (u *user) peerDisconnected(p *peer, err error) {
	if atomic.AddUintptr(&p.done, 1) > 1 {
		return
	}

	// TODO penalize the client ip, within memory

	if p.url != "" {
		u.db.PeerDisconnected(p.url, err)
	}
}

func (u *user) peerDownlink(p *peer) {
	defer p.Close()

	for {
		var any anypb.Any
		err := recvMessage(p.reader, &any)
		if err != nil {
			u.peerDisconnected(p, err)
			return
		}

		n, err := any.UnmarshalNew()
		if err != nil {
			u.peerDisconnected(p, err)
			return
		}

		switch msg := n.(type) {
		case *pb.RequestMsg:
		case *pb.RequestPeer:
			if !validatePoW(msg.Token, msg.Pow) {
				u.peerDisconnected(p, fmt.Errorf("invalid pow"))
				return
			}
			url := u.db.GetByToken(msg.Token)
			go p.SendPeer(url)
		case *pb.MsgList:
		case *pb.MsgContainer:
			// 1. validate hash
			hash := hasher(msg.Msg)
			if !bytes.Equal(hash[:], msg.MsgHash) {
				u.peerDisconnected(p, fmt.Errorf("invalid message hash"))
				return
			}
			// 2. validate pow
			if !validatePoW(msg.MsgHash, msg.Pow) {
				u.peerDisconnected(p, fmt.Errorf("invalid pow"))
				return
			}
			// 3. try decode
			var m pb.Message
			err := proto.Unmarshal(msg.Msg, &m)
			if err != nil {
				u.peerDisconnected(p, err)
				return
			}
			// 4. store and try decrypt
			if m.TargetCohort == u.cohort {
				u.db.StoreMessage(msg, m.TargetCohort)
				rMsg := u.DecryptMessage(&m)
				if rMsg != nil {
					u.db.StoreDecryptedMessage(rMsg)
				}
			} else if rand.Float64() < epsilon {
				u.db.StoreMessage(msg, m.TargetCohort)
			}
		case *pb.ResponsePeer:
			_ = msg
		default:
			u.peerDisconnected(p, fmt.Errorf("unknown message type %T", msg))
			return
		}
	}
}

func (u *user) peerUplink(p *peer) {
	defer p.Close()

	listTimer := time.NewTicker(u.cfg.ListMessageTime)
	defer listTimer.Stop()

	for {
		var msg proto.Message
		select {
		case msg = <-p.queue:
		case <-listTimer.C:
			continue // TODO
		case <-p.ctx.Done():
			break
		}

		any, err := anypb.New(msg)
		if err != nil {
			panic(err)
		}

		err = sendMessage(p.writer, any)
		if err != nil {
			u.peerDisconnected(p, err)
			return
		}
	}
}

func (u *user) runPeer(p *peer) *peer {
	p.queue = make(chan proto.Message)
	go u.peerUplink(p)
	go u.peerDownlink(p)
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

	return u.runPeer(&peer{
		ctx:    ctx,
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

	return u.runPeer(&peer{
		ctx:    ctx,
		reader: r,
		writer: w,
		cohort: ok.Cohort,
		tokens: ok.PeerTokens,
		key:    sKey,
	}), nil
}
