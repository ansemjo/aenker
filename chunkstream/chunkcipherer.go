// Package chunkstream provides a Writer and a Reader to encrypt or decrypt
// data in small authenticated chunks.
//
// You probably want ae.NewWriter() and ae.NewReader() rather than these chunkStreamer
// wrappers. The AEAD is used with a simple nonce counter, so you MUST to provide a
// unique key if you use this package directly.
package chunkstream

import (
	"crypto/cipher"

	"golang.org/x/crypto/chacha20poly1305"
)

// NewAEAD is the function to instantiate a new authenticated encryption cipher to be used by the
// chunkStreamer. If you want to use a different one, assign one before calling either
// NewReader or NewWriter, e.g.:
//  AEAD = func(key []byte) (aead cipher.AEAD, err error) {
//  	block, err := aes.NewCipher(key)
//  	if err != nil {
//  		return
//  	}
//  	return cipher.NewGCM(block)
//  }
var NewAEAD func([]byte) (cipher.AEAD, error) = chacha20poly1305.New

func init() {
}

// chunkCipherer is the cryptographic core of a chunked Reader or Writer.
type chunkCipherer struct {
	cipher cipher.AEAD
	ctr    *nonceCounter
	info   []byte
}

// NewChunkCipherer instantiates a new AEAD cipher and returns it in a ChunkCipherer
// struct, together with associated data that will be used for every chunk and a NonceCounter
// that will be incremented on every call to Seal() or Open().
// A ChunkCipherer should only ever be used to only seal or only open, as both functions
// share the same NonceCounter.
func newChunkCipherer(key, info []byte) (*chunkCipherer, error) {

	cc := &chunkCipherer{info: info}
	var err error

	cc.cipher, err = NewAEAD(key)
	if err != nil {
		return nil, err
	}

	cc.ctr = newNonceCounter(cc.cipher.NonceSize())
	return cc, err

}

func (cc *chunkCipherer) Seal(plain []byte) (ciphertext []byte) {
	return cc.cipher.Seal(plain[:0], cc.ctr.Next(), plain, cc.info)
}

func (cc *chunkCipherer) Open(ciphertext []byte) (plaintext []byte, err error) {
	return cc.cipher.Open(ciphertext[:0], cc.ctr.Next(), ciphertext, cc.info)
}
