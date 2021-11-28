package nymo

import (
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
)

const (
	nymoName     = "nymo"
	nymoVersion  = 1  // v0.1
	blockSize    = 32 // AES-256
	cohortNumber = 64
	bitStrength  = sha256.Size*8 - 20
	protoPrefix  = nymoName + "://"
	epsilon  = 0.1 // 10% messages/peers will be out-of-cohort
)

var (
	curve    = elliptic.P256() // NIST P-256
	cReader  = rand.Reader
	hasher   = sha256.Sum256
	encoding = binary.BigEndian

	curveByteLen = (curve.Params().BitSize + 7) / 8 // 256
)
