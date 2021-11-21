package nymo

import (
	"encoding/binary"
	"io"
	"unsafe"

	"google.golang.org/protobuf/proto"
)

const uint32Size = int(unsafe.Sizeof(uint32(0)))

var encoding = binary.BigEndian

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
