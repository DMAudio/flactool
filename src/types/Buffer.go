package types

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"sync"
)

type Buffer struct {
	buffer *bytes.Buffer
	locker sync.RWMutex
}

func NewBuffer() *Buffer {
	return &Buffer{
		buffer: bytes.NewBuffer([]byte{}),
	}
}

func (b *Buffer) Dump() ([]byte, error) {
	if b == nil {
		return nil,fmt.Errorf("不能导出未初始化的Buffer对象")
	}
	return ioutil.ReadAll(b.buffer)
}

func (b *Buffer) Read(p []byte) (n int, err error) {
	b.locker.RLock()
	defer b.locker.RUnlock()
	return b.buffer.Read(p)
}

func (b *Buffer) Write(p []byte) (n int, err error) {
	b.locker.Lock()
	defer b.locker.Unlock()
	return b.buffer.Write(p)
}

func (b *Buffer) WriteString(str string) (n int, err error) {
	b.locker.Lock()
	defer b.locker.Unlock()
	return b.buffer.WriteString(str)
}

func (b *Buffer) String() string {
	b.locker.RLock()
	defer b.locker.RUnlock()
	return b.buffer.String()
}

func (b *Buffer) WriteStringsIE(parts ...string) {
	b.locker.Lock()
	defer b.locker.Unlock()
	for _, part := range parts {
		b.buffer.WriteString(part)
	}
}

func (b *Buffer) WriteStrings(parts ...string) (n int, err error) {
	b.locker.Lock()
	defer b.locker.Unlock()

	sum := 0
	for _, part := range parts {
		n, err := b.buffer.WriteString(part)
		sum += n
		if err != nil {
			return sum, err
		}
	}
	return sum, nil
}
