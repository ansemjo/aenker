package aenker

import ce "github.com/ansemjo/aenker/error"

const (

	// ChunkSize is the amount of plaintext that is encrypted per chunk
	ChunkSize = 256

	// ErrExtraData indicates that there is extra data appended to the ciphertext
	ErrExtraData = ce.ConstError("aenker: extraneous data after ciphertext")
)
