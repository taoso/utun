package utun

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestXor(t *testing.T) {
	cases := []struct {
		data   []byte
		key    []byte
		result []byte
	}{
		{
			data: []byte{0x11, 0x22, 0x33, 0x44, 0x55},
			key:  []byte{0xaa, 0xbb},
			result: []byte{
				0x11 ^ 0xaa,
				0x22 ^ 0xbb,
				0x33 ^ 0xaa,
				0x44 ^ 0xbb,
				0x55 ^ 0xaa,
			},
		},
		{
			data:   []byte{},
			key:    []byte{0xaa, 0xbb},
			result: []byte{},
		},
		{
			data:   []byte{0x11},
			key:    []byte{},
			result: []byte{0x11},
		},
	}

	for i, c := range cases {
		xor(c.data, c.key)

		if !bytes.Equal(c.data, c.result) {
			t.Fatal("invalid result", i, c.data)
		}
	}
}

func TestXor2(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)

	data := make([]byte, 256)
	rand.Read(data)

	data2 := make([]byte, 256)
	copy(data2, data)

	xor(data, key)
	xor2(data, key)

	if !bytes.Equal(data, data2) {
		t.Fatal("invalid result", data)
	}
}

func BenchmarkXor(b *testing.B) {
	key := make([]byte, 32)
	rand.Read(key)

	data := make([]byte, 256)
	rand.Read(data)

	for range b.N {
		xor(data, key)
	}
}

func BenchmarkXor2(b *testing.B) {
	key := make([]byte, 32)
	rand.Read(key)

	data := make([]byte, 256)
	rand.Read(data)

	for range b.N {
		xor2(data, key)
	}
}
