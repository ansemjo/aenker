package aenker

import (
	"crypto/cipher"

	"golang.org/x/crypto/chacha20poly1305"
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
