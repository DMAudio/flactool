package types

import (
	"bufio"
	"io"
	"io/ioutil"
	"strconv"
)

type BinaryReader struct {
	rs io.Reader
}

var TMBinary_CanNotReadBytes = NewMask(
	"CANNOT_READ_BYTES",
	"无法读取定长数据（长度：{{length}}）",
)
var TMBinary_CanNotSkipBytes = NewMask(
	"CANNOT_Skip_BYTES",
	"无法跳过定长数据（长度：{{length}}）",
)


func NewBinaryReader(rs io.Reader) *BinaryReader {
	br := BinaryReader{
		rs: rs,
	}

	return &br
}

func (r *BinaryReader) ReadBytes(length uint64) ([]byte, *Exception) {
	bytes := make([]byte, length)
	_, err := r.rs.Read(bytes)
	if err != nil {
		return nil, NewException(TMBinary_CanNotReadBytes, map[string]string{
			"length": strconv.FormatUint(length, 10),
		}, err)
	}

	return bytes, nil
}

func (r *BinaryReader) Skip(length int) *Exception {
	br := bufio.NewReader(r.rs)
	if _, err := br.Discard(length); err != nil {
		return NewException(TMBinary_CanNotSkipBytes, map[string]string{
			"length": strconv.Itoa(length),
		}, err)
	}

	return nil
}

func (r *BinaryReader) ReadAllFollowedBytes() ([]byte, *Exception) {
	if output, err := ioutil.ReadAll(r.rs); err != nil {
		return nil, NewException(TMBinary_CanNotReadBytes, map[string]string{
			"length": "余下全部",
		}, err)
	} else {
		return output, nil
	}

}
