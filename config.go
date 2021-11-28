package nymo

import (
	"log"
	"time"
)

type Config struct {
	Logger          *log.Logger
	MaxPeerToken    int
	ListMessageTime time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		Logger:          log.Default(),
		MaxPeerToken:    10,
		ListMessageTime: time.Minute * 5,
	}
}
