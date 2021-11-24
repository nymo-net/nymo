package nymo

import (
	"crypto/ecdsa"
	"crypto/x509"
	"fmt"
)

type user struct {
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

func (u *user) Export() ([]byte, error) {
	return x509.MarshalECPrivateKey(u.key)
}

func ImportUser(der []byte) (*user, error) {
	key, err := x509.ParseECPrivateKey(der)
	if err != nil {
		return nil, err
	}
	if key.Curve != curve {
		return nil, fmt.Errorf("pkey not using %s", curve.Params().Name)
	}
	return &user{cohort: getCohort(key.X, key.Y), key: key}, nil
}

func GenerateUser() (*user, error) {
	key, err := ecdsa.GenerateKey(curve, cReader)
	if err != nil {
		return nil, err
	}
	return &user{cohort: getCohort(key.X, key.Y), key: key}, nil
}
