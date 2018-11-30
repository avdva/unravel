// Copyright 2018 Aleksandr Demakin. All rights reserved.

package pjw

import (
	"hash"
)

// PJW is an implementation of a non-cryptographic hash function
// created by Peter J. Weinberger.
// See
// https://en.wikipedia.org/wiki/PJW_hash_function
type PJW uint32

// New returns new PJW.
func New() hash.Hash32 {
	return new(PJW)
}

// Write appends hashed bytes to the existing hash.
func (pjw *PJW) Write(p []byte) (n int, err error) {
	const val = 0xF0000000
	h := uint32(*pjw)
	for _, c := range p {
		h = (h << 4) + uint32(c)
		high := h & val
		if high != 0 {
			h = h ^ (high >> 24)
		}
		h = h & (^high)
	}
	*pjw = PJW(h)
	return len(p), nil
}

// Sum appends hash to 'in' and returns the resulting slice.
func (pjw *PJW) Sum(in []byte) []byte {
	v := uint32(*pjw)
	return append(in, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

// Reset sets hash to zero.
func (pjw *PJW) Reset() {
	*pjw = 0
}

// Size returns hash size in bytes.
func (pjw PJW) Size() int {
	return 4
}

// BlockSize returns block size.
func (pjw *PJW) BlockSize() int {
	return 1
}

// Sum32 returns hash as a uint32 value.
func (pjw *PJW) Sum32() uint32 {
	return uint32(*pjw)
}
