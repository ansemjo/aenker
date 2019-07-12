// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

// +build pbkdf

package keyderivation

import (
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/blake2b"
)

// Password derives a 32 byte key from a password and salt with Argon2i and
// the predefined cost settings time=32, memory=256MB, threads=4.
// Keys generated this way are compatible with https://github.com/ansemjo/stdkdf.
func Password(password []byte, salt string) (key []byte) {
	s := blake2b.Sum256([]byte(salt))
	return argon2.Key(password, s[:], 32, 256*1024, 4, 32)
}
