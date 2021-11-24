package nymo

import (
	"bytes"
	"encoding/binary"
	"math/big"
)

func calcPoW(data []byte) uint64 {
	dataLen := len(data)
	powData := make([]byte, dataLen+8)
	copy(powData, data)

	var counter uint64
	var bInt big.Int
	for counter < 0xFFFFFFFF {
		encoding.PutUint64(powData[dataLen:], counter)
		hash := hasher(powData)
		bInt.SetBytes(hash[:])
		if bInt.BitLen() < bitStrength {
			break
		}
		counter++
	}

	return counter
}

func validatePoW(data []byte, pow uint64) bool {
	buf := bytes.NewBuffer(data)
	_ = binary.Write(buf, encoding, pow)
	var bInt big.Int
	hash := hasher(buf.Bytes())
	bInt.SetBytes(hash[:])
	return bInt.BitLen() < bitStrength
}
