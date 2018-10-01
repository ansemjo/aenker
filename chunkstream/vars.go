package chunkstream

import (
	"golang.org/x/crypto/chacha20poly1305"
)

// AEAD is the authenticated encryption cipher used by the chunkStreamer
var AEAD = chacha20poly1305.New
