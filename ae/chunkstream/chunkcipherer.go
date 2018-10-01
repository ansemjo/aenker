package chunkstream

import "crypto/cipher"

// chunkCipherer is the cryptographic core of a chunked reader or writer
type chunkCipherer struct {
	cipher cipher.AEAD
	ctr    *nonceCounter
	info   []byte
}

func newChunkCipherer(key, info []byte) (*chunkCipherer, error) {

	cc := &chunkCipherer{info: info}
	var err error

	cc.cipher, err = AEAD(key)
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
