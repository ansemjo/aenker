// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package ae

import (
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/hkdf"
	"hash"
	"io"
)

func KDF(userkey, salt []byte, info string) (key []byte) {

	// lambda function to conform to hkdf's hash type
	blake := func() hash.Hash {
		h, err := blake2b.New384(nil)
		if err != nil {
			// since we don't use a key, this shouldn't happen
			panic(err)
		}
		return h
	}

	hkdf := hkdf.New(blake, key, salt, []byte(info))

	key = make([]byte, 32)
	_, err := io.ReadFull(hkdf, key)
	if err != nil {
		// we only need 32 bytes .. i'm confident that there
		// should be no error
		panic(err)
	}

	return

}
