package nymo

import (
	"crypto/ecdsa"
	"crypto/tls"
	"math/big"
	"sync"
	"time"
)

type user struct {
	cfg    Config
	db     Database
	cohort uint32
	key    *ecdsa.PrivateKey
	cert   tls.Certificate
	id     [hashTruncate]byte

	peerLock sync.RWMutex
	peers    map[[hashTruncate]byte]*peer
	total    uint
	numIn    uint
	retry    peerRetrier
}

func (u *user) Address() *address {
	return &address{
		cohort: u.cohort,
		x:      u.key.X,
		y:      u.key.Y,
	}
}

func (u *user) Run() {
	for {
		u.dialNewPeers()
		time.Sleep(u.cfg.ScanPeerTime)
	}
}

func OpenUser(db Database, userKey []byte, cert tls.Certificate, cfg *Config) *user {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	key := new(ecdsa.PrivateKey)
	key.Curve = curve
	key.D = new(big.Int).SetBytes(userKey)
	key.X, key.Y = curve.ScalarBaseMult(userKey)

	hash := hasher(cert.Certificate[0])
	return &user{
		cfg:    *cfg,
		db:     db,
		cohort: getCohort(key.X, key.Y),
		key:    key,
		cert:   cert,
		id:     truncateHash(hash[:]),
		peers:  make(map[[hashTruncate]byte]*peer),
		retry:  peerRetrier{m: make(map[string]time.Time)},
	}
}

func GenerateUser() ([]byte, error) {
	key, err := ecdsa.GenerateKey(curve, cReader)
	if err != nil {
		return nil, err
	}
	return key.D.Bytes(), nil
}
