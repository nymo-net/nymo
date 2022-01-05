package nymo

import (
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"math"
	"time"
	"unsafe"
)

const (
	nymoName     = "nymo"
	nymoVersion  = 1  // v0.1
	blockSize    = 32 // AES-256
	digestSize   = sha256.Size
	cohortNumber = 64
	hashTruncate = 8
	bitStrength  = sha256.Size*8 - 22
	protoPrefix  = nymoName + "://"
	epsilon      = 0.1 // 10% messages/peers will be out-of-cohort

	msgListMax  = 500
	peerListMax = 20

	uint16Size    = int(unsafe.Sizeof(uint16(0)))
	maxPacketSize = math.MaxUint16 - uint16Size // 64 KiB
)

var (
	curve    = elliptic.P256() // NIST P-256
	cReader  = rand.Reader
	hasher   = sha256.Sum256
	encoding = binary.BigEndian

	curveByteLen = (curve.Params().BitSize + 7) / 8 // 256

	emptyTime time.Time
)
