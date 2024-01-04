package utun

import (
	"bytes"
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
