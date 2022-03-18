package nymo

import (
	"log"
	"time"
)

type Config struct {
	MaxInCohortConn  uint
	MaxOutCohortConn uint

	ListMessageTime time.Duration
	ScanPeerTime    time.Duration
	PeerRetryTime   time.Duration

	Logger *log.Logger

	LocalPeerAnnounce bool
	LocalPeerDiscover bool
	VerifyServerCert  bool
}

func DefaultConfig() *Config {
	return &Config{
		MaxInCohortConn:  5,
		MaxOutCohortConn: 5,

		ListMessageTime: time.Minute * 5,
		ScanPeerTime:    time.Second * 30,
		PeerRetryTime:   time.Minute,
	}
}
