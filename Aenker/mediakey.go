// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package aenker

import (
	"io"

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

// length of chunksize encoding
const chunkSizeEnc = 4

// length of a sealed MEK blob
const blobSize = nonceSize + keySize + chunkSizeEnc + aeadOverhead

// NewKey returns a new random key for usage with Aenker
func NewKey() (key []byte) {
	return randomBytes(keySize)
}

// NewMEK generates a new random media encryption key (MEK)
// and seals it with the supplied key encryption key (KEK). It
// returns the plain MEK and a blob which is needed later
// for decryption.
func (ae *Aenker) sealNewMEK(w io.Writer) (lengthWritten int, err error) {

	mek := NewKey()                 // generate a random media encryption key
	ae.mek = &mek                   // save mek in struct
	nonce := randomBytes(nonceSize) // generate a random nonce

	aead, err := ae.cipher.NewX(*ae.kek) // init AEAD, use the 'X' variant with larger nonce
	if err != nil {                      // size so random nonces can safely be used
		return
	}

	plain := append(*ae.mek, itob(uint32(ae.chunksize))...) // concatenate MEK || chunksize
	ad := []byte(mediaEncryptionAD)                         // get associated data
	sealed := aead.Seal(nil, nonce, plain, ad)              // encrypt the concatenation
	blob := append(nonce, sealed...)                        // concatenate nonce and sealed MEK
	lengthWritten, err = w.Write(blob)                      // write blob to writer

	return

}

// OpenMEK decrypts a previously sealed MEK blob and returns
// the plain MEK within.
func (ae *Aenker) openMEK(r io.Reader) (err error) {

	aead, err := ae.cipher.NewX(*ae.kek) // init AEAD, use the 'X' variant with larger nonce
	if err != nil {                      // size so random nonces can safely be used
		return
	}

	ad := []byte(mediaEncryptionAD) // get associated data

	nonce := make([]byte, nonceSize) // allocate slice for nonce
	_, err = io.ReadFull(r, nonce)   // read nonce from reader
	if err != nil {
		return
	}

	sealed := make([]byte, blobSize-nonceSize) // allocate slice for the sealed MEK with auth tag
	_, err = io.ReadFull(r, sealed)            // read sealed data from reader
	if err != nil {
		return
	}

	plain, err := aead.Open(nil, nonce, sealed, ad) // decrypt media encryption key
	if err != nil {                                 // bail if decryption failed
		return
	}

	mek := plain[:keySize] // get MEK from plaintext slice
	ae.mek = &mek          // save MEK in struct

	ae.chunksize = int(btoi(plain[keySize : keySize+chunkSizeEnc])) // decode chunksize

	return

}
