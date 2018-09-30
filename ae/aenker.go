// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package ae

import (
	"crypto/cipher"
	"io"

	"github.com/ansemjo/aenker/ae/chunkstream"
	"golang.org/x/crypto/chacha20poly1305"
)

const (
	magic            = "aenker"
	kdfInfoSymmetric = magic + " symmetric"
	kdfInfoPassword  = magic + " password"
	kdfInfoElliptic  = magic + " elliptic"
)

type aeadCipher struct {
	New  func([]byte) (cipher.AEAD, error)
	NewX func([]byte) (cipher.AEAD, error)
}

// Aenker is a struct that supports a sort of "streamed" AEAD usage,
// where the plaintext is split into equal parts and encrypted with
// an interMAClib-like construction
type Aenker struct {
	cipher    aeadCipher
	kek       *[]byte
	mek       *[]byte
	chunksize int
}

// NewAenker returns a new Aenker, ready for
// encryption or decryption with the given key
func NewAenker(key *[]byte) (ae *Aenker) {
	aead := aeadCipher{
		New:  chacha20poly1305.New,
		NewX: chacha20poly1305.NewX,
	}
	return &Aenker{aead, key, nil, -1}
}

// Encrypt encrypt data on reader and writes ciphertext to writer.
// It generates a new random media encryption key and nonce, seals the
// MEK with the key given in NewAenker(), splits the data in the given reader into
// chunks and encrypts them individually with the MEK.
// Data written to writer is 'nonce || sealedMEKBlob || sealedChunks[]' and must be
// passed entirely to Decrypt() later for successful decryption.
func (ae *Aenker) Encrypt(w io.Writer, r io.Reader, chunksize int) (lengthWritten uint64, err error) {

	ae.chunksize = chunksize

	mek, err := ae.sealNewMEK(w)
	if err != nil {
		return
	}

	chunks, err := chunkstream.Encrypt(chunkstream.Options{
		Key:       *ae.mek,
		Info:      nil,
		Reader:    r,
		Writer:    w,
		Chunksize: chunksize,
	})
	lengthWritten = uint64(mek) + chunks
	return

}

// Decrypt reads ciphertext from reader and writes plaintext to writer.
// It reads the nonce and sealedMEKBlob from the reader, attempts to decrypt
// the MEK with the key given to NewAenker() and then reads the sealed chunks
// from the reader and decrypts them individually.
// Any truncation, bit-flips or chunk reordering is detected as an authentication
// error. Additional data after an authenticated ciphertext is an ErrExtraData warning.
func (ae *Aenker) Decrypt(w io.Writer, r io.Reader) (lengthWritten uint64, err error) {

	err = ae.openMEK(r)
	if err != nil {
		return
	}

	return chunkstream.Decrypt(chunkstream.Options{
		Key:       *ae.mek,
		Info:      nil,
		Reader:    r,
		Writer:    w,
		Chunksize: ae.chunksize,
	})

}
