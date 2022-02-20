package nymo

import (
	"log"
	"time"
)

type Config struct {
	MaxConcurrentConn uint
	ListMessageTime   time.Duration
	ScanPeerTime      time.Duration
	PeerRetryTime     time.Duration
	Logger            *log.Logger

	LocalPeerAnnounce bool
	LocalPeerDiscover bool
}

func DefaultConfig() *Config {
	return &Config{
		MaxConcurrentConn: 10,
		ListMessageTime:   time.Minute * 5,
		ScanPeerTime:      time.Second * 30,
		PeerRetryTime:     time.Minute,
	}
}
