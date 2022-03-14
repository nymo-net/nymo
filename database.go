package nymo

import "github.com/nymo-net/nymo/pb"

type PeerHandle interface {
	AddKnownMessages([]*pb.Digest) []*pb.Digest
	ListMessages(size uint) []*pb.Digest
	AckMessages()
	AddKnownPeers([]*pb.Digest) []*pb.Digest
	ListPeers(size uint) []*pb.Digest
	Disconnect(error)
}

type PeerEnumerate interface {
	Url() string
	Cohort() uint32
	Next(error) bool
	Connect(id []byte, cohort uint32) PeerHandle
	Close()
}

type Database interface {
	ClientHandle(id []byte) PeerHandle
	AddPeer(url string, digest *pb.Digest)
	EnumeratePeers() PeerEnumerate
	GetUrlByHash(urlHash []byte) (url string)

	GetMessage(hash []byte) (msg []byte, pow uint64)
	IgnoreMessage(digest *pb.Digest)
	StoreMessage(hash []byte, c *pb.MsgContainer, f func() (cohort uint32, err error)) error
	StoreDecryptedMessage(*Message)
}
