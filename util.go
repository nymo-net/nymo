package nymo

import (
	"bytes"
	"encoding/binary"
	"io"
	"net/http"
	"sync"
	"time"
	"unsafe"

	"google.golang.org/protobuf/proto"
)

func sameCohort(target, pov uint32) bool {
	return pov == target || target == cohortNumber
}

func (u *User) peerSameCohort(peer uint32) bool {
	return sameCohort(peer, u.cohort)
}

func sendMessage(conn io.Writer, m proto.Message) error {
	data, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	if len(data) > maxPacketSize {
		panic("exceed size")
	}

	buf := make([]byte, uint16Size+len(data))
	encoding.PutUint16(buf, uint16(len(data)))
	copy(buf[uint16Size:], data)

	_, err = conn.Write(buf)
	return err
}

func recvMessage(conn io.Reader, m proto.Message) error {
	var size uint16
	err := binary.Read(conn, encoding, &size)
	if err != nil {
		return err
	}

	buf := make([]byte, size)
	_, err = io.ReadFull(conn, buf)
	if err != nil {
		return err
	}

	return proto.Unmarshal(buf, m)
}

type writeFlusher struct {
	w io.Writer
	f http.Flusher
}

func (w *writeFlusher) Write(p []byte) (n int, err error) {
	defer w.f.Flush()
	return w.w.Write(p)
}

func padBlock(input []byte) []byte {
	pad := (len(input)/blockSize+1)*blockSize - len(input)
	return append(input, bytes.Repeat([]byte{byte(pad)}, pad)...)
}

func trimBlock(input []byte) []byte {
	c := input[len(input)-1]
	b := int(c)
	if b <= 0 || b > blockSize {
		return nil
	}
	for i := len(input) - b; i < len(input)-1; i++ {
		if input[i] != c {
			return nil
		}
	}
	return input[:len(input)-b]
}

func truncateHash(hash []byte) [hashTruncate]byte {
	return *(*[hashTruncate]byte)(unsafe.Pointer(&hash[0]))
}

type peerRetrier struct {
	l sync.Mutex
	m map[string]time.Time
}

func (p *peerRetrier) addSelf(url string) {
	p.l.Lock()
	defer p.l.Unlock()

	p.m[url] = emptyTime
}

func (p *peerRetrier) add(url string, timeout time.Duration) {
	ddl := time.Now().Add(timeout)

	p.l.Lock()
	defer p.l.Unlock()

	t, ok := p.m[url]
	if !ok || (t != emptyTime && t.Before(ddl)) {
		p.m[url] = ddl
	}
}

func (p *peerRetrier) noRetry(url string) bool {
	p.l.Lock()
	defer p.l.Unlock()

	t, ok := p.m[url]
	if !ok {
		return false
	}
	if t == emptyTime || time.Until(t) > 0 {
		return true
	}
	delete(p.m, url)
	return false
}
