package keyderivation

import (
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/curve25519"
)

// Symmetric derives a key from a high-entropy secret and a salt
func Symmetric(secret, salt []byte) (key []byte) {

	return Derive(secret, salt, "symmetric")

}

// Password derives a key from a password and salt with argon2i
func Password(password, salt []byte) (key []byte) {

	pwhash := argon2.Key(password, salt, 16, 64*1024, 2, 32)
	return Derive(pwhash, nil, "password")

}

// Elliptic derives a key by performing diffie-hellman with my private and a peer's public key over curve25519
func Elliptic(private, peer []byte) (key []byte) {

	// get a fixed 32 byte array from slice
	get32 := func(slice []byte) (array *[32]byte) {
		array = new([32]byte)
		copy(array[:], slice)
		return
	}

	shared := new([32]byte)
	curve25519.ScalarMult(shared, get32(private), get32(peer))
	return Derive(shared[:], nil, "elliptic")

}
