package nymo

import (
	"bytes"
	"crypto/tls"
	"math/rand"
	"testing"

	"github.com/nymo-net/nymo/pb"
	"google.golang.org/protobuf/proto"
)

func TestMsgEncDec(t *testing.T) {
	key, err := GenerateUser()
	if err != nil {
		t.Fatal(err)
	}
	var db testDB
	u := OpenUser(&db, key, tls.Certificate{Certificate: [][]byte{{}}}, nil)

	buf := make([]byte, 2048)
	rand.Read(buf)
	err = u.NewMessage(u.Address(), buf)
	if err != nil {
		t.Fatal(err)
	}

	message := u.decryptMessage(db.last)
	if message == nil || !bytes.Equal(message.Content, buf) {
		t.Fatal("decrypt error")
	}
}

func TestMsgEncNoDec(t *testing.T) {
	key, err := GenerateUser()
	if err != nil {
		t.Fatal(err)
	}
	key2, err := GenerateUser()
	if err != nil {
		t.Fatal(err)
	}
	var db testDB
	u := OpenUser(&db, key, tls.Certificate{Certificate: [][]byte{{}}}, nil)
	u2 := OpenUser(new(testDB), key2, tls.Certificate{Certificate: [][]byte{{}}}, nil)

	buf := make([]byte, 2048)
	rand.Read(buf)
	err = u.NewMessage(u.Address(), buf)
	if err != nil {
		t.Fatal(err)
	}

	message := u2.decryptMessage(db.last)
	if message != nil {
		t.Fatal("decrypt success")
	}
}

type testDB struct {
	last *pb.Message
}

func (t *testDB) ClientHandle(id [hashTruncate]byte) PeerHandle {
	panic("implement me")
}

func (t *testDB) AddPeer(url string, digest *pb.Digest) {
	panic("implement me")
}

func (t *testDB) EnumeratePeers() PeerEnumerate {
	panic("implement me")
}

func (t *testDB) GetUrlByHash(urlHash [hashTruncate]byte) (url string) {
	panic("implement me")
}

func (t *testDB) GetMessage(hash [hashSize]byte) (msg []byte, pow uint64) {
	panic("implement me")
}

func (t *testDB) IgnoreMessage(digest *pb.Digest) {
	panic("implement me")
}

func (t *testDB) StoreMessage(hash [hashSize]byte, c *pb.MsgContainer, f func() (cohort uint32, err error)) error {
	_, err := f()
	if err == nil {
		t.last = new(pb.Message)
		err = proto.Unmarshal(c.Msg, t.last)
	}
	return err
}

func (t *testDB) StoreDecryptedMessage(message *Message) {
	panic("implement me")
}
