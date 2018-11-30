// Copyright 2018 Aleksandr Demakin. All rights reserved.

package hash

import "github.com/avdva/unravel/hash/pjw"

// MakeHasher returns a funcion, that computes hash of a string.
// It returns nil, if there's no such algorithm.
func MakeHasher(alg string) func(string) []byte {
	switch alg {
	case "pjw":
		return func(s string) []byte {
			h := pjw.New()
			h.Write([]byte(s))
			return h.Sum(nil)
		}
	default:
		return nil
	}
}
