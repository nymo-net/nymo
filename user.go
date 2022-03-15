package nymo

import (
	"context"
	"crypto/ecdsa"
	"crypto/tls"
	"math/big"
	"sync"
	"time"

	"github.com/nymo-net/nymo/pb"
)

type User struct {
	cfg    Config
	db     Database
	cohort uint32
	key    *ecdsa.PrivateKey
	cert   tls.Certificate

	peerLock sync.RWMutex
	peers    map[[hashTruncate]byte]*peer
	total    uint
	numIn    uint
	retry    peerRetrier
}

func (u *User) Address() *Address {
	return &Address{
		cohort: u.cohort,
		x:      u.key.X,
		y:      u.key.Y,
	}
}

func (u *User) Run(ctx context.Context) {
	if u.cfg.LocalPeerAnnounce {
		go func() {
			if e := u.ipv4PeerAnnounce(ctx); e != nil {
				u.cfg.Logger.Print(e)
			}
		}()
		go func() {
			if e := u.ipv6PeerAnnounce(ctx); e != nil {
				u.cfg.Logger.Print(e)
			}
		}()
	}
	if u.cfg.LocalPeerDiscover {
		go func() {
			if e := u.ipv4PeerDiscover(ctx); e != nil {
				u.cfg.Logger.Print(e)
			}
		}()
		go func() {
			if e := u.ipv6PeerDiscover(ctx); e != nil {
				u.cfg.Logger.Print(e)
			}
		}()
	}

	for ctx.Err() == nil {
		u.dialNewPeers(ctx)
		t := time.NewTimer(u.cfg.ScanPeerTime)
		select {
		case <-t.C:
		case <-ctx.Done():
			t.Stop()
			return
		}
	}
}

func (u *User) AddPeer(url string) {
	hash := hasher([]byte(url))
	u.db.AddPeer(url, &pb.Digest{
		Hash:   hash[:hashTruncate],
		Cohort: cohortNumber, // XXX: when unknown, as wildcard
	})
}

func OpenSupernode(db Database, cert tls.Certificate, cfg *Config) *User {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	hash := hasher(cert.Certificate[0])
	return &User{
		cfg:    *cfg,
		db:     db,
		cohort: cohortNumber,
		cert:   cert,
		peers:  map[[hashTruncate]byte]*peer{truncateHash(hash[:]): nil},
		retry:  peerRetrier{m: make(map[string]time.Time)},
	}
}

func OpenUser(db Database, userKey []byte, cert tls.Certificate, cfg *Config) *User {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	key := new(ecdsa.PrivateKey)
	key.Curve = curve
	key.D = new(big.Int).SetBytes(userKey)
	key.X, key.Y = curve.ScalarBaseMult(userKey)

	hash := hasher(cert.Certificate[0])
	return &User{
		cfg:    *cfg,
		db:     db,
		cohort: getCohort(key.X, key.Y),
		key:    key,
		cert:   cert,
		peers:  map[[hashTruncate]byte]*peer{truncateHash(hash[:]): nil},
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
