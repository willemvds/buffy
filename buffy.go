package buffy

import (
	"errors"
	"io"
	"sync/atomic"
)

const defaultBufferSize = 64 * 1024

var EmptySlice []byte = []byte{}
var ErrNoBytes = errors.New("no bytes available at specified offset")

// Buffy is an append only byte buffer for a very specific use case where:
// 1) The bytes are never changed once written.
// 2) There is 1 writer.
// 3) There are multiple concurrent readers.
// This also means the buffer is never re-used.
// The api does not prevent callers from modifying the underlying bytes.
type Buffy struct {
	buffer atomic.Pointer[[]byte]
	closed atomic.Bool
}

func NewWithSize(initialSize int) *Buffy {
	if initialSize < 0 {
		initialSize = defaultBufferSize
	}
	buf := make([]byte, 0, defaultBufferSize)
	bfy := Buffy{}
	bfy.buffer.Store(&buf)

	return &bfy
}

func New() *Buffy {
	return NewWithSize(defaultBufferSize)
}

func (bfy *Buffy) IsClosed() bool {
	return bfy.closed.Load()
}

// Even though this can't error (only panic) we implement the io.Writer interface
// so it can be used idiomatically.
func (bfy *Buffy) Write(in []byte) (int, error) {
	newBuf := append(*bfy.buffer.Load(), in...)
	bfy.buffer.Swap(&newBuf)

	return len(in), nil
}

func (bfy *Buffy) Close() error {
	bfy.closed.Store(true)

	return nil
}

// Returns the slice where the data is stored. The caller should treat this slice
// as immutable and not modify it. If modification is needed the caller should first
// make their own local copy to work with.
func (bfy *Buffy) Bytes() []byte {
	return *bfy.buffer.Load()
}

func (bfy *Buffy) Since(n int) ([]byte, error) {
	buf := *bfy.buffer.Load()

	if bfy.IsClosed() {
		remainder := EmptySlice
		if n < len(buf) {
			remainder = buf[n:]
		}

		return remainder, io.EOF
	}

	if n >= len(buf) {
		return EmptySlice, ErrNoBytes
	}

	return buf[n:], nil
}
