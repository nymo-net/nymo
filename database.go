package nymo

type DatabaseFactory = func(userKey []byte) (Database, error)

type Database interface {
	GetUserKey() ([]byte, error)
	Close() error
}
