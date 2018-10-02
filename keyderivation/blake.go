// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package keyderivation

import (
	"hash"

	"golang.org/x/crypto/blake2b"
)

// translate blake2's keyed interfaces to simple unkeyed ones
func unkeyed(newkeyed func([]byte) (hash.Hash, error)) func() hash.Hash {

	return func() hash.Hash {
		h, err := newkeyed(nil)
		if err != nil {
			// since we don't use a key, this shouldn't happen
			panic(err)
		}
		return h
	}

}

// Unkeyed Blake2b Hashes for use in the crypto.hkdf.New constructor
var (
	Blake2b256 = unkeyed(blake2b.New256)
	Blake2b384 = unkeyed(blake2b.New384)
	Blake2b512 = unkeyed(blake2b.New512)
)
