package nymo

import "github.com/nymo-net/nymo/pb"

// PeerHandle is a helper object for Nymo core to handle a peer connection.
// Only one PeerHandle will be requested to any peer (with peer ID) before Disconnect is called.
//
// For the following documentation, 𝓜 and 𝓟 is the set of messages and peers stored in the database already,
// 𝓢 is the set of messages we know that the peer knows, and 𝓡 is the set of peers we know that the peer knows.
// For a newly encountered peer ID, 𝓢 = 𝓡 = ∅.
//
// Notice 𝓢 and 𝓡 requires consistency during the lifetime of a PeerHandle object, but isn't
// required to be consistent across restarts, as messages/peers lists should be re-listed if not found.
// However, for best performance, it is best if the implementation can persist 𝓢 and 𝓡 across restarts.
//
// It is then up to the implementation to optimize storage space usage. One way is to remove
// information of a peer that we haven't connected in a long time.
type PeerHandle interface {
	// AddKnownMessages adds a new batch of known messages 𝓑 by the peer to the database.
	// It should update 𝓢 <- 𝓢 ∪ 𝓑 and return 𝓑 \ 𝓜.
	AddKnownMessages([]*pb.Digest) []*pb.Digest
	// ListMessages returns 𝓓 = 𝓜 \ 𝓢, with size no more than the size parameter.
	// It should also record 𝓓 (see AckMessages for more information).
	ListMessages(size uint) []*pb.Digest
	// AckMessages confirms the receptions of the last ListMessages. That means
	// the implementation should update 𝓢 <- 𝓢 ∪ 𝓓.
	//
	// It is guaranteed that AckMessages will always be called between two ListMessages.
	AckMessages()
	// AddKnownPeers adds a new batch of known peers 𝓒 by the peer to the database.
	// It should update 𝓡 <- 𝓡 ∪ 𝓒 and return 𝓒 \ 𝓟.
	AddKnownPeers([]*pb.Digest) []*pb.Digest
	// ListPeers returns 𝓟 \ 𝓡, with size no more than the size parameter.
	ListPeers(size uint) []*pb.Digest
	// Disconnect is called when the peer disconnects, the implementation can release resources.
	//
	// The error argument is the error (if any) why disconnection happened.
	// The implementation can optionally use this information to rank and ban peers.
	Disconnect(error)
}

// PeerEnumerate is a helper object for Nymo core to iterate through
// known peer URLs. The order is not specified but the implementation
// can choose to list good-ranking peers in front. See Next for more information.
type PeerEnumerate interface {
	// Url returns the peer URL that is being currently iterated over.
	Url() string
	// Cohort returns the cached cohort of the peer URL that is being currently iterated over.
	Cohort() uint32
	// Next iterates to the next peer URL, and should return false
	// if it no longer has more peer URLs to iterate over.
	//
	// The error argument is the error (if any) happened trying to
	// connect to the peer URL that is being currently iterated over.
	// The implementation can optionally use this information to rank and ban peers.
	//
	// Next is guaranteed to be called once before the iteration with error being nil,
	// meaning Url and Cohort should return the first result available after a call to Next.
	Next(error) bool
	// Connect confirms the connection to the URL that is being currently iterated over.
	// The implementation should return a PeerHandle of the given peer ID, and update the cached
	// cohort corresponding the URL, if necessary.
	//
	// A call to Connect shouldn't automatically advance the iterator.
	Connect(id [hashTruncate]byte, cohort uint32) PeerHandle
	// Close stops the iteration, and the implementation should release resources.
	Close()
}

// Database needs to be implemented by any frontend for user/database
// related actions. Data encryption can be optionally implemented,
// transparent to the Nymo core.
//
// Each peer is uniquely identified by its peer ID. A peer can have multiple
// URLs, which is uniquely identified by either the URL string or the hash of it.
// However, there is no relationship between the URL and the peer ID.
//
// Each message is uniquely identified by its hash. However, the implementation
// should always use (hash, cohort) as the key as there might be adversaries
// trying to cheat on what cohort the message was sent to. The actual cohort is only
// confirmed when the message is received and stored by StoreMessage.
type Database interface {
	// ClientHandle returns a PeerHandle of the given peer ID.
	// See PeerHandle for more information.
	ClientHandle(id [hashTruncate]byte) PeerHandle
	// AddPeer adds a peer (URL, URL hash, cohort) information to the database.
	// Notice that the cohort is only a cache field for better performance
	// (see PeerEnumerate for more information).
	AddPeer(url string, digest *pb.Digest)
	// EnumeratePeers returns a PeerEnumerate of all the known peers (added by AddPeer).
	// See PeerEnumerate for more information.
	EnumeratePeers() PeerEnumerate
	// GetUrlByHash returns the URL by the given URL hash.
	// The implementation can return anything if the hash is not found in the database
	// (it should not error out).
	GetUrlByHash(urlHash [hashTruncate]byte) (url string)

	// GetMessage returns the message data and PoW of the given message hash (stored by StoreMessage).
	// The implementation can return anything if the hash is not found in the database
	// (it should not error out).
	GetMessage(hash [hashSize]byte) (msg []byte, pow uint64)
	// IgnoreMessage tells the database to ignore the existence of the message (message hash, cohort).
	// (i.e. take message as if it were already received and removed, see StoreMessage)
	IgnoreMessage(digest *pb.Digest)
	// StoreMessage stores the message (message data, message hash, PoW, cohort) into the database.
	//
	// Nymo core might store the same message twice into the database in case of a race.
	// The implementation can check first if the message with the hash is already in the database,
	// and choose not to do nothing (return nil directly). Otherwise, it should call function f.
	//
	// If f returns a non-nil error, StoreMessage should abort and return that error. Otherwise,
	// it should proceed and store the message into the database.
	//
	// The implementation can choose to remove ancient messages by still storing the hash but removing
	// the actual data.
	StoreMessage(hash [hashSize]byte, c *pb.MsgContainer, f func() (cohort uint32, err error)) error
	// StoreDecryptedMessage stores a decrypted message (see Message for more information) into the database.
	// It might be called synchronously when the f function is called by the implementation.
	//
	// The implementation should not block the call for long.
	StoreDecryptedMessage(*Message)
}
