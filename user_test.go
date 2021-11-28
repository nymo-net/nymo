package nymo

import (
	"testing"
)

func TestSend(t *testing.T) {
	u1, err := GenerateUser(newMemDb, nil)
	if err != nil {
		t.Error(err)
	}

	u2, err := GenerateUser(newMemDb, nil)
	if err != nil {
		t.Error(err)
	}

	if u1.key.Equal(u2.key) {
		t.Error("same key")
	}

	const s = "hello user2"
	msg, err := u1.NewMessage(u2.Address(), s)
	if err != nil {
		t.Error(err)
	}

	rMsg := u2.DecryptMessage(msg)
	if rMsg == nil || rMsg.Content != s {
		t.Error("message mismatch")
	}
}

func TestReopen(t *testing.T) {
	u, err := GenerateUser(newMemDb, nil)
	if err != nil {
		t.Error(err)
	}

	u2, err := OpenUser(u.db, nil)
	if err != nil {
		t.Error(err)
	}

	if u.cohort != u2.cohort || !u.key.Equal(u2.key) {
		t.Error("user mismatch")
	}
}

type memDb struct {
	NopDatabase
	key []byte
}

func (m *memDb) GetUserKey() ([]byte, error) {
	return m.key, nil
}

func newMemDb(userKey []byte) (Database, error) {
	return &memDb{key: userKey}, nil
}
