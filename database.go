package nymo

import "github.com/nymo-net/nymo/pb"

type DatabaseFactory = func(userKey []byte) (Database, error)

type Database interface {
	GetUserKey() ([]byte, error)

	AddPeer(url string, token []byte, cohort uint32, connected bool)
	GetByToken(token []byte) (url string)
	PeerDisconnected(url string, reason error)
	GetStoredPeers(cohort uint32, size uint) (tokens [][]byte, err error)

	StoreMessage(msg *pb.MsgContainer, cohort uint32)
	StoreDecryptedMessage(*Message)
	ListMessages(known [][]byte) ([]*pb.MsgDigest, error)
}

type NopDatabase struct{}

func (NopDatabase) GetUserKey() ([]byte, error) { panic("implement me") }

func (NopDatabase) AddPeer(string, []byte, uint32, bool) {}

func (NopDatabase) GetByToken([]byte) (url string) { return "" }

func (NopDatabase) PeerDisconnected(string, error) {}

func (NopDatabase) GetStoredPeers(uint32, uint) ([][]byte, error) { return nil, nil }

func (NopDatabase) StoreMessage(*pb.MsgContainer, uint32) {}

func (NopDatabase) StoreDecryptedMessage(*Message) {}

func (NopDatabase) ListMessages([][]byte) ([]*pb.MsgDigest, error) { return nil, nil }
