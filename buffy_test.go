package buffy

import (
	"bytes"
	"testing"
)

func TestEmptyBufferSince0(t *testing.T) {
	bfy := New()

	buf, err := bfy.Since(0)

	if !bytes.Equal(buf, EmptySlice) || err != ErrNoBytes {
		t.Errorf("Since(0) EXPECT ([]byte{}, %s),  GOT (%x, %s)", ErrNoBytes, buf, err)
	}
}

func TestSince(t *testing.T) {
	bfy := New()
	in := []byte("ABCDEF")
	bfy.Write(in)

	buf, err := bfy.Since(0)

	if !bytes.Equal(buf, in) || err != nil {
		t.Errorf("Since(0) EXPECT (%x, %s), GOT (%x, %s)", in, "", buf, err)
	}
}

func TestSinceIndex(t *testing.T) {
	bfy := New()
	in := []byte("ABCDEF")
	idx := 3
	expect := in[idx:]
	bfy.Write(in)

	buf, err := bfy.Since(idx)

	if !bytes.Equal(buf, expect) || err != nil {
		t.Errorf("Since(%d) EXPECT (%s, %s), GOT (%s, %s)", idx, expect, "", buf, err)
	}
}

func TestSincePastIndex(t *testing.T) {
	bfy := New()
	in := []byte("ABCDEF")
	idx := 6
	expect := EmptySlice
	expectErr := ErrNoBytes
	bfy.Write(in)

	buf, err := bfy.Since(idx)

	if !bytes.Equal(buf, expect) || err != expectErr {
		t.Errorf("Since(%d) EXPECT (%s, %s), GOT (%s, %s)", idx, expect, expectErr, buf, err)
	}
}
