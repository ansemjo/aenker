// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

// Package keyderivation provides functions to derive keys for symmetric encryption.
package keyderivation

import (
	"hash"
	"io"

	"golang.org/x/crypto/hkdf"
)

// Hash is the hash function used in HKDF. If you want a different one assign it before
// calling HKDF or any function that uses it. I.e.:
//  keyderivation.Hash = sha512.New
var Hash func() hash.Hash = Blake2b512

// HKDF wraps crypto/hkdf to generate a single 32 byte key with the Hash defined at package level.
func HKDF(secret, salt []byte, info string) (key []byte) {

	// instantiate hkdf
	hkdf := hkdf.New(Hash, secret, salt, []byte(info))

	// read 32 bytes
	key = make([]byte, 32)
	_, err := io.ReadFull(hkdf, key)
	if err != nil {
		// we only read 32 bytes ..
		// i'm fairly confident there should be no error
		panic(err)
	}

	return

}
