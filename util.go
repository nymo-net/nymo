package nymo

import (
	"bytes"
	"encoding/binary"
	"io"
	"net/http"
	"unsafe"

	"google.golang.org/protobuf/proto"
)

const uint32Size = int(unsafe.Sizeof(uint32(0)))

func sendMessage(conn io.Writer, m proto.Message) error {
	data, err := proto.Marshal(m)
	if err != nil {
		return err
	}

	buf := make([]byte, len(data)+uint32Size)
	encoding.PutUint32(buf, uint32(len(data)))
	copy(buf[uint32Size:], data)

	_, err = conn.Write(buf)
	return err
}

func recvMessage(conn io.Reader, m proto.Message) error {
	var size uint32
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
	b := int(input[len(input)-1])
	if b <= 0 || b > blockSize {
		return nil
	}
	return input[:len(input)-b]
}
