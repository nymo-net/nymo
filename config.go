package nymo

import (
	"time"
)

type Config struct {
	MaxPeerToken    int
	ListMessageTime time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		MaxPeerToken:    10,
		ListMessageTime: time.Minute * 5,
	}
}
