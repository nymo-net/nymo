package nymo

import "bytes"

// TODO

func calcPoW(data []byte) []byte {
	return data
}

func validatePoW(data, pow []byte) bool {
	return bytes.Equal(data, pow)
}
