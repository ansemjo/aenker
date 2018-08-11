package aenker

import (
	"golang.org/x/crypto/chacha20poly1305"
)

// associated data for the MEK encryption
const mediaEncryptionAD = "Aenker Media Encryption Key"

// length of the required nonce required for XChaCha20Poly1305 during MEK encryption
const nonceSize = chacha20poly1305.NonceSizeX

// length of the required key in bytes
const keySize = chacha20poly1305.KeySize

// length of added overhead by aead
const aeadOverhead = 16 // chacha.New().Overhead()

// MekBlobSize is the length of a sealed MEK blob
const MekBlobSize = nonceSize + keySize + aeadOverhead

// NewKey returns a new random key for usage with Aenker
func NewKey() (key []byte) {
	return randomBytes(keySize)
}

// NewMEK generates a new random media encryption key (MEK)
// and seals it with the supplied key encryption key (KEK). It
// returns the plain MEK and a blob which is needed later
// for decryption.
// TODO: add chunksize to blob
func NewMEK(kek []byte) (mek, blob []byte, err error) {

	mek = NewKey()                  // generate a random media encryption key
	nonce := randomBytes(nonceSize) // generate a random nonce

	aead, err := chacha20poly1305.NewX(kek) // init AEAD, use the 'X' variant with larger nonce
	if err != nil {                         // size so random nonces can safely be used
		return
	}

	ad := []byte(mediaEncryptionAD)          // associated data
	sealed := aead.Seal(nil, nonce, mek, ad) // encrypt the MEK
	blob = append(nonce, sealed...)          // concatenate nonce and sealed MEK
	return

}

// OpenMEK decrypts a previously sealed MEK blob and returns
// the plain MEK within.
func OpenMEK(kek, blob []byte) (mek []byte, err error) {

	aead, err := chacha20poly1305.NewX(kek) // init AEAD, use the 'X' variant with larger nonce
	if err != nil {                         // size so random nonces can safely be used
		return
	}

	ad := []byte(mediaEncryptionAD)              // associated data
	nonce := blob[:nonceSize]                    // nonce is in first part
	sealed := blob[nonceSize:]                   // the rest is the sealed MEK with auth tag
	mek, err = aead.Open(nil, nonce, sealed, ad) // decrypt media encryption key
	return

}
