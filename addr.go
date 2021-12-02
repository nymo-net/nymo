package nymo

import (
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/base64"
	"math/big"
	"strings"
)

type address struct {
	cohort uint32
	x, y   *big.Int
}

func (r *address) Cohort() uint32 {
	return r.cohort
}

func (r *address) Bytes() []byte {
	return elliptic.MarshalCompressed(curve, r.x, r.y)
}

func (r *address) String() string {
	return protoPrefix + base64.RawURLEncoding.EncodeToString(r.Bytes())
}

func getCohort(x, y *big.Int) uint32 {
	hash := sha256.New()
	hash.Write(x.Bytes())
	hash.Write(y.Bytes())

	var h big.Int
	h.SetBytes(hash.Sum(nil))
	h.Mod(&h, big.NewInt(cohortNumber))

	return uint32(h.Uint64()) + 1
}

func NewAddress(addr string) *address {
	if !strings.HasPrefix(addr, protoPrefix) {
		return nil
	}

	buf, err := base64.RawURLEncoding.DecodeString(addr[len(protoPrefix):])
	if err != nil {
		return nil
	}

	x, y := elliptic.UnmarshalCompressed(curve, buf)
	if x == nil {
		return nil
	}
	return newAddress(x, y)
}

func newAddress(x, y *big.Int) *address {
	return &address{getCohort(x, y), x, y}
}
