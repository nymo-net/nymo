package nymo

import (
	"math/rand"
	"testing"
)

func TestPoW(t *testing.T) {
	buf := make([]byte, 32)
	_, err := rand.Read(buf)
	if err != nil {
		t.Error(err)
	}

	w := calcPoW(buf)
	if !validatePoW(buf, w) {
		t.Error("invalid pow")
	}
}

func TestPoWRandom(t *testing.T) {
	buf := make([]byte, 32)
	_, err := rand.Read(buf)
	if err != nil {
		t.Error(err)
	}

	if validatePoW(buf, rand.Uint64()) {
		t.Error("invalid pow")
	}
}