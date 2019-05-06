package util

func BytesToUInt64(bytes []byte) uint64 {
	var result uint64
	for i := 0; i < IntMin(8, len(bytes)); i++ {
		result = result | uint64(bytes[len(bytes)-i-1])<<uint(8*i)
	}

	return result
}

func BytesToInt64(bytes []byte) uint64 {
	var result uint64
	for i := 0; i < IntMin(8, len(bytes)); i++ {
		result = result | uint64(bytes[len(bytes)-i-1])<<uint(7*i)
	}

	return result
}