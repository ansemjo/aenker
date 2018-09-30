// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package keyderivation

import (
	"io"

	"golang.org/x/crypto/hkdf"
)

// Hash is the hash function used by Derive
var Hash = Blake2b384

// Info is the customization prefix that is used with the derivation type (Info+" $type")
var Info = "aenker"

// Derive uses HKDF with the given Hash to generate a 32 byte key
func Derive(secret, salt []byte, infosuffix string) (key []byte) {

	// instantiate hkdf
	hkdf := hkdf.New(Hash, secret, salt, []byte(Info+" "+infosuffix))

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
