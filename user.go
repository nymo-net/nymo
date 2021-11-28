package nymo

import (
	"crypto/ecdsa"
	"math/big"
)

type user struct {
	cfg    Config
	db     Database
	cohort uint32
	key    *ecdsa.PrivateKey
}

func (u *user) Address() *address {
	return &address{
		cohort: u.cohort,
		x:      u.key.X,
		y:      u.key.Y,
	}
}

func openUser(db Database, key *ecdsa.PrivateKey, cfg *Config) *user {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	return &user{
		cfg:    *cfg,
		db:     db,
		cohort: getCohort(key.X, key.Y),
		key:    key,
	}
}

func OpenUser(db Database, cfg *Config) (*user, error) {
	d, err := db.GetUserKey()
	if err != nil {
		return nil, err
	}

	key := new(ecdsa.PrivateKey)
	key.Curve = curve
	key.D = new(big.Int).SetBytes(d)
	key.X, key.Y = curve.ScalarBaseMult(d)

	return openUser(db, key, cfg), nil
}

func GenerateUser(factory DatabaseFactory, cfg *Config) (*user, error) {
	key, err := ecdsa.GenerateKey(curve, cReader)
	if err != nil {
		return nil, err
	}

	db, err := factory(key.D.Bytes())
	if err != nil {
		return nil, err
	}

	return openUser(db, key, cfg), nil
}
