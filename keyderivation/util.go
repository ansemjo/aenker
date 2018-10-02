package keyderivation

import "golang.org/x/crypto/curve25519"

// // get a fixed 32 byte array from slice
// func get32(slice []byte) (array *[32]byte) {
// 	array = new([32]byte)
// 	copy(array[:], slice)
// 	return
// }

// Public returns the Curve25519 public key of a secret key.
func Public(sec *[32]byte) (pub *[32]byte) {
	pub = new([32]byte)
	curve25519.ScalarBaseMult(pub, sec)
	return pub
}
