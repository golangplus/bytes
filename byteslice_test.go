// Copyright 2015 The Golang Plus Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bytesp

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"testing"
	"unicode/utf8"

	"github.com/golangplus/testing/assert"
)

func ExampleByteSlice() {
	var b ByteSlice
	b.WriteByte(65)
	b.WriteString("bc")

	fmt.Println(b)
	fmt.Println(string(b))
	// OUTPUT:
	// [65 98 99]
	// Abc
}

func TestByteSlice(t *testing.T) {
	var bs ByteSlice
	assert.Equal(t, "len(bs)", len(bs), 0)
	assert.StringEqual(t, "bs", bs, "[]")

	bs.Write([]byte{1, 2, 3})
	assert.Equal(t, "len(bs)", len(bs), 3)
	assert.StringEqual(t, "bs", bs, "[1 2 3]")

	p := make([]byte, 2)
	bs.Read(p)
	assert.Equal(t, "len(bs)", len(bs), 1)
	assert.StringEqual(t, "bs", bs, "[3]")
	assert.StringEqual(t, "p", p, "[1 2]")

	bs.Read(make([]byte, 1))
	assert.Equal(t, "len(bs)", len(bs), 0)
	assert.StringEqual(t, "bs", bs, "[]")

	bs.Write([]byte{4, 5})
	assert.Equal(t, "len(bs)", len(bs), 2)
	assert.StringEqual(t, "bs", bs, "[4 5]")

	bs.WriteByte(6)

	c, err := bs.ReadByte()
	assert.Equal(t, "c", c, byte(4))
	assert.Equal(t, "err", err, nil)
	assert.StringEqual(t, "bs", bs, "[5 6]")

	bs.WriteRune('A')
	assert.Equal(t, "len(bs)", len(bs), 3)
	assert.StringEqual(t, "bs", bs, "[5 6 65]")
	bs.WriteRune('中')
	assert.Equal(t, "len(bs)", len(bs), 6)
	assert.StringEqual(t, "bs", bs, "[5 6 65 228 184 173]")

	bs.WriteString("世界")
	assert.Equal(t, "len(bs)", len(bs), 12)
	assert.StringEqual(t, "bs", bs, "[5 6 65 228 184 173 228 184 150 231 149 140]")

	bs.Skip(1)
	assert.StringEqual(t, "bs", bs, "[6 65 228 184 173 228 184 150 231 149 140]")

	bs.Close()

	bs = nil
	fmt.Fprint(&bs, "ABC")
	assert.StringEqual(t, "bs", bs, "[65 66 67]")

	data := make([]byte, 35*1024)
	io.ReadFull(rand.Reader, data)
	bs = nil
	n, err := bs.ReadFrom(bytes.NewReader(data))
	assert.Equal(t, "err", err, nil)
	assert.Equal(t, "n", n, int64(len(data)))
	assert.Equal(t, "bs == data", bytes.Equal(bs, data), true)

	bs = nil
	n, err = ByteSlice(data).WriteTo(&bs)
	assert.Equal(t, "err", err, nil)
	assert.Equal(t, "n", n, int64(len(data)))
	assert.Equal(t, "bs == data", bytes.Equal(bs, data), true)

	bs = []byte("A中文")
	r, size, err := bs.ReadRune()
	assert.Equal(t, "err", err, nil)
	assert.Equal(t, "size", size, 1)
	assert.Equal(t, "r", r, 'A')
	r, size, err = bs.ReadRune()
	assert.Equal(t, "err", err, nil)
	assert.Equal(t, "size", size, len([]byte("中")))
	assert.Equal(t, "r", r, '中')
	r, size, err = bs.ReadRune()
	assert.Equal(t, "err", err, nil)
	assert.Equal(t, "size", size, len([]byte("文")))
	assert.Equal(t, "r", r, '文')
}

func TestByteSlice_Bug_Read(t *testing.T) {
	var s ByteSlice
	n, err := s.Read(make([]byte, 1))
	t.Logf("n: %d, err: %v", n, err)
	assert.Equal(t, "n", 0, 0)
	assert.Equal(t, "err", err, io.EOF)
}

func TestByteSlice_Bug_ReadRune(t *testing.T) {
	s := ByteSlice{65, 0xff, 66}
	r, sz, err := s.ReadRune()
	assert.Equal(t, "r", r, 'A')
	assert.Equal(t, "sz", sz, 1)
	assert.Equal(t, "err", err, nil)
	r, sz, err = s.ReadRune()
	assert.Equal(t, "r", r, utf8.RuneError)
	assert.Equal(t, "sz", sz, 1)
	assert.Equal(t, "err", err, nil)

	r, sz, err = s.ReadRune()
	assert.Equal(t, "r", r, 'B')
	assert.Equal(t, "sz", sz, 1)
	assert.Equal(t, "err", err, nil)
}

func TestByteSlice_WriteItoa(t *testing.T) {
	var s ByteSlice
	s.WriteItoa(1234, 10)
	s.WriteItoa(255, 16)

	assert.Equal(t, "s", string(s), "1234ff")
}

func BenchmarkByteSlice100(b *testing.B) {
	var data [100]byte
	for i := 0; i < b.N; i++ {
		b := ByteSlice(data[:])
		for {
			if _, err := b.ReadByte(); err != nil {
				break
			}
		}
	}
}

func BenchmarkBuffer100(b *testing.B) {
	var data [100]byte
	for i := 0; i < b.N; i++ {
		b := bytes.NewBuffer(data[:])
		for {
			if _, err := b.ReadByte(); err != nil {
				break
			}
		}
	}
}

func BenchmarkByteSlice10(b *testing.B) {
	var data [10]byte
	for i := 0; i < b.N; i++ {
		b := ByteSlice(data[:])
		for {
			if _, err := b.ReadByte(); err != nil {
				break
			}
		}
	}
}

func BenchmarkBuffer10(b *testing.B) {
	var data [10]byte
	for i := 0; i < b.N; i++ {
		b := bytes.NewBuffer(data[:])
		for {
			if _, err := b.ReadByte(); err != nil {
				break
			}
		}
	}
}
