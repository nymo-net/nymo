package nymo

import (
	"time"
)

type Config struct {
	MaxConcurrentConn uint
	MaxPeerToken      int
	ListMessageTime   time.Duration
	ScanPeerTime      time.Duration
	PeerRetryTime     time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		MaxConcurrentConn: 10,
		MaxPeerToken:      10,
		ListMessageTime:   time.Minute * 5,
		ScanPeerTime:      time.Second * 30,
		PeerRetryTime:     time.Minute,
	}
}
