// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package keyderivation

import (
	"io"

	"golang.org/x/crypto/hkdf"
)

// Derive uses HKDF to generate a 32 byte key. Info will be
// assembled as keyderivation.Info+" "+info.
func Derive(secret, salt []byte, info string) (key []byte) {

	// instantiate hkdf
	hkdf := hkdf.New(Hash, secret, salt, []byte(Info+" "+info))

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
