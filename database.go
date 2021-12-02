package nymo

import "github.com/nymo-net/nymo/pb"

type PeerHandle interface {
	AddKnownMessages([]*pb.Digest) (need []*pb.Digest, ignored [][]byte)
	AckKnownMessages([][]byte)
	ListMessages(size uint) []*pb.Digest
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
	MessageStat(cohort uint32) (in uint, out uint)

	ClientHandle(id []byte) PeerHandle
	AddPeer(url string, digest *pb.Digest)
	EnumeratePeers() PeerEnumerate
	GetUrlByHash(urlHash []byte) (url string)

	GetMessage(hash []byte) (msg []byte, pow uint64)
	StoreMessage(hash []byte, c *pb.MsgContainer, f func() (cohort uint32, err error)) error
	StoreDecryptedMessage(*Message)
}
