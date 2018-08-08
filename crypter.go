package main

import (
	"bufio"
	"crypto/cipher"
	"io"

	"golang.org/x/crypto/chacha20poly1305"
)

const (
	CHUNKSIZE = 27
)

type crypter struct {
	aead cipher.AEAD
}

// init an aead crypter
func newCrypter() *crypter {
	// init with zero key
	key := make([]byte, chacha20poly1305.KeySize)
	aead, err := chacha20poly1305.New(key)
	if err != nil {
		panic("AEAD init failed")
	}
	return &crypter{aead: aead}
}

func isLast(r *bufio.Reader) bool {
	_, err := r.Peek(1)
	return err != nil
}

func (c *crypter) Encrypt(w io.Writer, r io.Reader) (nOut int64, errOut error) {

	// buffered reader
	br := bufio.NewReader(r)

	// internal chunk size and counter
	chunk := make([]byte, CHUNKSIZE-1)
	nonce := NewNonceCounter()

	for {

		n, err := io.ReadFull(br, chunk[:CHUNKSIZE-1])
		last := isLast(br)
		if n > 0 {
			if last {
				stderr("this is the last chunk")
			}
			stderr(sfmt("next chunk %d, %d bytes", nonce.ctr, n))
			chunk = Pad(chunk[:n], CHUNKSIZE, last)
			ct := c.aead.Seal(nil, nonce.Next(), chunk[:n], nil)
			nw, err := w.Write(ct)
			nOut += int64(nw)
			if err != nil {
				errOut = err
				return
			}
		}
		if false { //err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		} else if err != nil {
			errOut = err
			return
		}

	}
	return

}
