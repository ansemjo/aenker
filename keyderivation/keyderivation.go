package keyderivation

import (
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/curve25519"
)

// Symmetric derives a 32 byte key from a high-entropy secret and a salt.
func Symmetric(secret, salt []byte) (key []byte) {
	return Derive(secret, salt, "symmetric")
}

// Password derives a 32 byte key from a password and salt with Argon2i and
// the predefined cost settings time=16, memory=64MB, threads=2.
func Password(password, salt []byte) (key []byte) {
	pwhash := argon2.Key(password, salt, 16, 64*1024, 2, 32)
	return Derive(pwhash, nil, "password")
}

// Elliptic derives a key by performing Diffie-Hellman with a private and a
// peer's public key over Curve25519. Will use [32]byte arrays internally, so
// make sure to pass 32 byte slices and especially do not use nil for peer!
func Elliptic(private, peer []byte) (key []byte) {
	shared := new([32]byte)
	curve25519.ScalarMult(shared, get32(private), get32(peer))
	return Derive(shared[:], nil, "elliptic")
}
