package aenker

import (
	ce "github.com/ansemjo/aenker/error"
	"golang.org/x/crypto/chacha20poly1305"
)

const (

	// KeyLength is the length of the required key in bytes
	KeyLength = chacha20poly1305.KeySize

	// ChunkSize is the amount of plaintext that is encrypted per chunk
	//ChunkSize = 16 * 1024

	// ErrExtraData indicates that there is extra data appended to the ciphertext
	ErrExtraData = ce.ConstError("aenker: extraneous data after ciphertext")
)
