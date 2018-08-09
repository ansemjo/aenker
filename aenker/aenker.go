package aenker

import (
	"crypto/cipher"

	constError "github.com/ansemjo/aenker/error"
	"golang.org/x/crypto/chacha20poly1305"
)

const (

	// ChunkSize is the amount of plaintext that is encrypted per chunk
	ChunkSize = 256

	// ErrExtraData indicates that there is extra data appended to the ciphertext
	ErrExtraData = constError.Error("aenker: extraneous data after ciphertext")
)

// Aenker is a struct that supports a sort of "streamed" AEAD usage,
// where the plaintext is split into equal parts and encrypted with
// an interMAClib-like construction
type Aenker struct {
	aead cipher.AEAD
}

// NewAenker return a new Aenker with initialized cipher
func NewAenker(key []byte) (ae *Aenker, keyerror error) {
	aead, err := chacha20poly1305.New(key)
	return &Aenker{aead}, err
}
