package types

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func IntMin(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func UIntToBytes(n uint64, b byte) ([]byte, error) {
	var err error
	bytesBuffer := bytes.NewBuffer(make([]byte, b))
	switch b {
	case 1:
		tmp := uint8(n)
		err = binary.Write(bytesBuffer, binary.BigEndian, &tmp)
	case 2:
		tmp := uint16(n)
		err = binary.Write(bytesBuffer, binary.BigEndian, &tmp)
	case 3, 4:
		tmp := uint32(n)
		err = binary.Write(bytesBuffer, binary.BigEndian, &tmp)
	case 5, 6, 7, 8:
		tmp := uint64(n)
		err = binary.Write(bytesBuffer, binary.BigEndian, &tmp)
	}

	if err == nil {
		oBytes := bytesBuffer.Bytes()
		return oBytes[len(oBytes)-int(b):], nil
	}
	return nil, fmt.Errorf("UIntToBytes b param is invaild")
}
