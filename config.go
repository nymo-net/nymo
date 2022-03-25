package nymo

import (
	"log"
	"time"
)

// Config statically configures several parameters of Nymo core.
type Config struct {
	// MaxInCohortConn is the number of maximum number of
	// possible in-cohort connections.
	MaxInCohortConn uint
	// MaxOutCohortConn is the number of maximum number of
	// possible out-of-cohort connections.
	MaxOutCohortConn uint

	// ListMessageTime is the interval at which messages are listed to the peers.
	ListMessageTime time.Duration
	// ScanPeerTime is the interval at which new peer connections are tried.
	ScanPeerTime time.Duration
	// PeerRetryTime is the interval at which a peer connection will be retried.
	PeerRetryTime time.Duration

	// Logger is a custom control over logging outputs.
	Logger *log.Logger

	// whether Local Peer Announce is enabled.
	LocalPeerAnnounce bool
	// whether Local Peer Discover is enabled.
	LocalPeerDiscover bool
	// whether server peer certificate should be validated (against its domain name).
	VerifyServerCert bool
}

// DefaultConfig returns a copy of default config for itemized modification.
func DefaultConfig() *Config {
	return &Config{
		MaxInCohortConn:  5,
		MaxOutCohortConn: 5,

		ListMessageTime: time.Minute * 5,
		ScanPeerTime:    time.Second * 30,
		PeerRetryTime:   time.Minute,
	}
}
