// Copyright 2018 Aleksandr Demakin. All rights reserved.

package pjw

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPJW(t *testing.T) {
	const (
		hello = "Hello"
		world = ", world!"
	)
	r := require.New(t)
	pjw := New()
	r.Equal(4, pjw.Size())
	r.Equal(1, pjw.BlockSize())
	r.Equal(uint32(0), pjw.Sum32())
	r.Equal([]byte{0, 0, 0, 0}, pjw.Sum(nil))

	pjw.Reset()

	written, err := pjw.Write([]byte(hello))
	r.NoError(err)
	r.Equal(len(hello), written)
	r.Equal(uint32(0x004ec32f), pjw.Sum32())
	r.Equal([]byte{0, 0x4e, 0xc3, 0x2f}, pjw.Sum(nil))

	written, err = pjw.Write([]byte(world))
	r.NoError(err)
	r.Equal(len(world), written)
	r.Equal(uint32(0x0925c3c1), pjw.Sum32())
	r.Equal([]byte{0x09, 0x25, 0xc3, 0xc1}, pjw.Sum(nil))

	pjw.Reset()

	r.Equal(uint32(0), pjw.Sum32())
	written, err = pjw.Write([]byte(hello + world))
	r.NoError(err)
	r.Equal(len(hello)+len(world), written)
	r.Equal(uint32(0x925c3c1), pjw.Sum32())
	r.Equal([]byte{0x01, 0x00, 0xff, 0x09, 0x25, 0xc3, 0xc1}, pjw.Sum([]byte{0x01, 0x00, 0xff}))
}
