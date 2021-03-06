package nymo

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/nymo-net/nymo/pb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type peer struct {
	ctx    context.Context
	handle PeerHandle

	reader io.ReadCloser
	writer io.Writer
	cohort uint32
	peers  []*pb.Digest
	key    []byte

	queue chan proto.Message

	msgReq  sync.Map
	msgQ    chan struct{}
	peerReq sync.Map

	msgProc uint32
}

func (p *peer) sendProto(msg proto.Message) {
	select {
	case <-p.ctx.Done():
	case p.queue <- msg:
	}
}

func (p *peer) requestMsg(diff []*pb.Digest, u *User) {
	count := 0
	for _, digest := range diff {
		// deal with out-of-cohort message
		if !sameCohort(u.cohort, digest.Cohort) {
			// should not receive if msg cohort is same as peer
			if sameCohort(p.cohort, digest.Cohort) {
				continue
			}

			// 1-ε chance of ignoring an out-of-cohort message
			if rand.Float64() >= epsilon {
				u.db.IgnoreMessage(digest)
				continue
			}
		}

		p.msgReq.Store(*(*[hashSize]byte)(unsafe.Pointer(&digest.Hash[0])), digest.Cohort)
		p.sendProto(&pb.RequestMsg{Hash: digest.Hash})
		count++
	}

	for count > 0 {
		if _, ok := <-p.msgQ; !ok {
			break
		}
		count--
	}
	atomic.StoreUint32(&p.msgProc, 0)
	p.sendProto(new(pb.MsgListAck))
}

func (p *peer) requestPeer(pred func(uint32) bool) bool {
	for i, digest := range p.peers {
		if !pred(digest.Cohort) {
			continue
		}
		p.peerReq.Store(truncateHash(digest.Hash), digest)
		p.peers = append(p.peers[:i], p.peers[i+1:]...)
		p.sendProto(&pb.RequestPeer{
			UrlHash: digest.Hash,
			Pow:     calcPoW(append(p.key, digest.Hash...)),
		})
		return true
	}
	return false
}

func (u *User) peerDownlink(p *peer) error {
	defer close(p.msgQ)
	defer p.reader.Close()

	var listTimer unsafe.Pointer

	defer func() {
		t := (*time.Timer)(atomic.SwapPointer(&listTimer, nil))
		if t != nil {
			t.Stop()
		}
	}()

	listMsg := func() {
		if atomic.SwapPointer(&listTimer, nil) != nil {
			p.sendProto(&pb.MsgList{Messages: p.handle.ListMessages(msgListMax)})
		}
	}

	for p.ctx.Err() == nil {
		var any anypb.Any
		err := recvMessage(p.reader, &any)
		if err != nil {
			return err
		}

		n, err := any.UnmarshalNew()
		if err != nil {
			return err
		}

		switch msg := n.(type) {
		case *pb.RequestMsg:
			if len(msg.Hash) != hashSize {
				return fmt.Errorf("unexpected hash length")
			}
			cont := new(pb.MsgContainer)
			cont.Msg, cont.Pow = u.db.GetMessage(copyHash(msg.Hash))
			p.sendProto(cont)
		case *pb.RequestPeer:
			if len(msg.UrlHash) != hashTruncate {
				return fmt.Errorf("unexpected hash length")
			}
			if !validatePoW(append(p.key, msg.UrlHash...), msg.Pow) {
				return fmt.Errorf("invalid pow")
			}
			p.sendProto(&pb.ResponsePeer{Address: u.db.GetUrlByHash(truncateHash(msg.UrlHash))})
		case *pb.ResponsePeer:
			hash := hasher([]byte(msg.Address))
			digest, loaded := p.peerReq.LoadAndDelete(truncateHash(hash[:]))
			if !loaded {
				return fmt.Errorf("unexpected peer response")
			}
			u.db.AddPeer(msg.Address, digest.(*pb.Digest))
		case *pb.PeerList:
			if p.peers != nil {
				return fmt.Errorf("unexpected peer list")
			}
			for _, l := range msg.Peers {
				if len(l.Hash) != hashTruncate || l.Cohort > cohortNumber {
					return fmt.Errorf("unexpected hash length")
				}
			}
			p.peers = p.handle.AddKnownPeers(msg.Peers)
		case *pb.MsgList:
			if atomic.LoadUint32(&p.msgProc) != 0 {
				return fmt.Errorf("unexpected msg list")
			}
			for _, l := range msg.Messages {
				if len(l.Hash) != hashSize || l.Cohort >= cohortNumber {
					return fmt.Errorf("unexpected hash length")
				}
			}
			p.msgProc = 1
			go p.requestMsg(p.handle.AddKnownMessages(msg.Messages), u)
		case *pb.MsgListAck:
			if atomic.LoadPointer(&listTimer) != nil {
				return errors.New("unexpected peer msg ack")
			}
			p.handle.AckMessages()
			atomic.StorePointer(&listTimer, unsafe.Pointer(time.AfterFunc(u.cfg.ListMessageTime, listMsg)))
		case *pb.MsgContainer:
			if msg.Msg == nil {
				return fmt.Errorf("no msg")
			}
			// 1. retrieve request
			msgHash := hasher(msg.Msg)
			expCohort, loaded := p.msgReq.LoadAndDelete(msgHash)
			if !loaded {
				return fmt.Errorf("unexpected msg response")
			}
			p.msgQ <- struct{}{}
			err := u.db.StoreMessage(msgHash, msg, func() (uint32, error) {
				// 2. validate pow
				if !validatePoW(msgHash[:], msg.Pow) {
					return 0, errors.New("invalid pow")
				}
				// 3. try decode
				var m pb.Message
				err := proto.Unmarshal(msg.Msg, &m)
				if err != nil {
					return 0, err
				}
				if m.TargetCohort != expCohort.(uint32) {
					return 0, errors.New("unexpected pow")
				}
				// 4. try decrypt and store
				if m.TargetCohort == u.cohort {
					if rMsg := u.decryptMessage(&m); rMsg != nil {
						u.db.StoreDecryptedMessage(rMsg)
					}
				}
				return m.TargetCohort, nil
			})
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown message type %T", msg)
		}
	}

	return p.ctx.Err()
}

func (u *User) peerUplink(p *peer) {
	defer p.reader.Close()

	for {
		var msg proto.Message
		select {
		case msg = <-p.queue:
		case <-p.ctx.Done():
			return
		}

		any, err := anypb.New(msg)
		if err != nil {
			panic(err)
		}

		err = sendMessage(p.writer, any)
		if err != nil {
			p.handle.Disconnect(err)
			return
		}
	}
}

func (u *User) runPeer(p *peer) *peer {
	p.queue = make(chan proto.Message, 10)
	p.queue <- &pb.PeerList{Peers: p.handle.ListPeers(peerListMax)}
	p.queue <- &pb.MsgList{Messages: p.handle.ListMessages(msgListMax)}
	go func() {
		p.handle.Disconnect(u.peerDownlink(p))
	}()
	go u.peerUplink(p)
	return p
}

func (u *User) newPeerAsServer(
	ctx context.Context, handle PeerHandle,
	r io.ReadCloser, w io.Writer,
	material []byte, cohort uint32) (*peer, error) {
	if err := sendMessage(w, &pb.HandshakeOK{
		Cohort: u.cohort,
	}); err != nil {
		return nil, err
	}

	return u.runPeer(&peer{
		ctx:    ctx,
		handle: handle,
		reader: r,
		writer: w,
		cohort: cohort,
		key:    material,
		msgQ:   make(chan struct{}, 100),
	}), nil
}

func (u *User) newPeerAsClient(
	ctx context.Context, handle PeerEnumerate,
	r io.ReadCloser, w io.Writer,
	id [hashTruncate]byte, sKey []byte) (*peer, error) {
	var ok pb.HandshakeOK
	if err := recvMessage(r, &ok); err != nil {
		return nil, err
	}

	return u.runPeer(&peer{
		ctx:    ctx,
		handle: handle.Connect(id, ok.Cohort),
		reader: r,
		writer: w,
		cohort: ok.Cohort,
		key:    sKey,
		msgQ:   make(chan struct{}, 100),
	}), nil
}
