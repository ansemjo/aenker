// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package keyderivation

import (
	"golang.org/x/crypto/curve25519"
)

// Elliptic perform anonymous Diffie-Hellman and then derives a 32 byte
// key from the resulting shared secret with HKDF. Salt and info may be nil but
// provide additional entropy for HKDF.
func Elliptic(private, peer *[32]byte, salt []byte, info string) (key []byte) {

	// perform anonymous diffie-hellman
	shared := new([32]byte)
	curve25519.ScalarMult(shared, private, peer)

	// derive key with hkdf
	return HKDF(shared[:], salt, info)

}

// Public returns the Curve25519 public key of a secret key.
func Public(sec *[32]byte) (pub *[32]byte) {
	pub = new([32]byte)
	curve25519.ScalarBaseMult(pub, sec)
	return pub
}
