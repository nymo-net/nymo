package nymo

import (
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/base64"
	"math/big"
	"strings"
)

// Address represents the address of a Nymo user. This struct should not be initialized directly.
type Address struct {
	cohort uint32
	x, y   *big.Int
}

// Cohort returns the cohort number of the address.
func (r *Address) Cohort() uint32 {
	return r.cohort
}

// Bytes returns the encoded address.
func (r *Address) Bytes() []byte {
	return elliptic.MarshalCompressed(curve, r.x, r.y)
}

// ConvertAddrToStr returns the string version of an encoded Nymo user address.
func ConvertAddrToStr(addr []byte) string {
	// first 6 bits is always 0, so truncate
	return protoPrefix + base64.RawURLEncoding.EncodeToString(addr)[1:]
}

// String returns the string version of the address.
// It is equivalent to ConvertAddrToStr(address.Bytes())
func (r *Address) String() string {
	return ConvertAddrToStr(r.Bytes())
}

func getCohort(x, y *big.Int) uint32 {
	hash := sha256.New()
	hash.Write(x.Bytes())
	hash.Write(y.Bytes())

	var h big.Int
	h.SetBytes(hash.Sum(nil))
	h.Mod(&h, big.NewInt(cohortNumber))

	return uint32(h.Uint64())
}

// NewAddress converts the string version of an address into the internal representation.
func NewAddress(addr string) *Address {
	if !strings.HasPrefix(addr, protoPrefix) {
		return nil
	}

	// first 6 bits should always be 0 ("A" in base 64)
	buf, err := base64.RawURLEncoding.DecodeString("A" + addr[len(protoPrefix):])
	if err != nil {
		return nil
	}
	return NewAddressFromBytes(buf)
}

// NewAddressFromBytes converts the encoded address into the internal representation.
func NewAddressFromBytes(addr []byte) *Address {
	x, y := elliptic.UnmarshalCompressed(curve, addr)
	if x == nil {
		return nil
	}
	return newAddress(x, y)
}

func newAddress(x, y *big.Int) *Address {
	return &Address{getCohort(x, y), x, y}
}
