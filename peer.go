package nymo

import (
	"context"
	"errors"
	"fmt"
	"io"
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
	peerReq sync.Map
}

func (p *peer) sendProto(msg proto.Message) {
	select {
	case <-p.ctx.Done():
	case p.queue <- msg:
	}
}

func (p *peer) requestMsg(diff []*pb.Digest, u *user) {
	in, out := u.db.MessageStat(u.cohort)

	// count in-cohort messages, and
	// remove all same-as-source cohort messages
	end := len(diff)
	for i := 0; i < end; i++ {
		if sameCohort(u.cohort, diff[i].Cohort) {
			in++
		} else if sameCohort(p.cohort, diff[i].Cohort) {
			end--
			diff[i], diff[end] = diff[end], diff[i]
			i--
		}
	}
	diff = diff[:end]

	quota := int(float64(in)/(1-epsilon)*epsilon) - int(out)
	for _, digest := range diff {
		if digest.Cohort != u.cohort {
			if quota <= 0 {
				continue
			}
			quota--
		}

		p.msgReq.Store(*(*[digestSize]byte)(unsafe.Pointer(&digest.Hash[0])), nil)
		p.sendProto(&pb.RequestMsg{Hash: digest.Hash})
	}
}

func (p *peer) requestPeer(cohort uint32) bool {
	if len(p.peers) > 0 {
		return false
	}
	for i, digest := range p.peers {
		if !sameCohort(digest.Cohort, cohort) {
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

func (u *user) peerDownlink(p *peer) error {
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
			cont := new(pb.MsgContainer)
			cont.Msg, cont.Pow = u.db.GetMessage(msg.Hash)
			p.sendProto(cont)
		case *pb.RequestPeer:
			if !validatePoW(append(p.key, msg.UrlHash...), msg.Pow) {
				return fmt.Errorf("invalid pow")
			}
			p.sendProto(&pb.ResponsePeer{Address: u.db.GetUrlByHash(msg.UrlHash)})
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
				if len(l.Hash) != hashTruncate {
					return fmt.Errorf("unexpected peer")
				}
			}
			p.peers = p.handle.AddKnownPeers(msg.Peers)
		case *pb.MsgList:
			diff, known := p.handle.AddKnownMessages(msg.Messages)
			p.sendProto(&pb.MsgKnown{Hashes: known})
			go p.requestMsg(diff, u)
		case *pb.MsgKnown:
			if atomic.LoadPointer(&listTimer) != nil {
				return errors.New("unexpected peer known msg")
			}
			p.handle.AckKnownMessages(msg.Hashes)
			atomic.StorePointer(&listTimer, unsafe.Pointer(time.AfterFunc(u.cfg.ListMessageTime, listMsg)))
		case *pb.MsgContainer:
			// 1. retrieve request
			msgHash := hasher(msg.Msg)
			_, loaded := p.msgReq.LoadAndDelete(msgHash)
			if !loaded {
				return fmt.Errorf("unexpected msg response")
			}
			err := u.db.StoreMessage(msgHash[:], msg, func() (uint32, error) {
				// 2. validate pow
				if !validatePoW(msgHash[:], msg.Pow) {
					return 0, fmt.Errorf("invalid pow")
				}
				// 3. try decode
				var m pb.Message
				err := proto.Unmarshal(msg.Msg, &m)
				if err != nil {
					return 0, err
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

func (u *user) peerUplink(p *peer) {
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

func (u *user) runPeer(p *peer) *peer {
	p.queue = make(chan proto.Message, 10)
	p.queue <- &pb.PeerList{Peers: p.handle.ListPeers(peerListMax)}
	p.queue <- &pb.MsgList{Messages: p.handle.ListMessages(msgListMax)}
	go func() {
		p.handle.Disconnect(u.peerDownlink(p))
	}()
	go u.peerUplink(p)
	return p
}

func (u *user) newPeerAsServer(
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
	}), nil
}

func (u *user) newPeerAsClient(
	ctx context.Context, handle PeerEnumerate,
	r io.ReadCloser, w io.Writer,
	id []byte, sKey []byte) (*peer, error) {
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
	}), nil
}
