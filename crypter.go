package main

import (
	"bufio"
	"crypto/cipher"
	"io"

	"golang.org/x/crypto/chacha20poly1305"
)

const (
	CHUNKSIZE = 24
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
	size := CHUNKSIZE
	chunk := make([]byte, size)
	nonce := NewNonceCounter()

	for {

		n, err := io.ReadFull(br, chunk[:size-1])
		last := isLast(br)
		if n > 0 {
			if last {
				//stderr("this is the last chunk")
			}
			stderr(sfmt("output chunk % 3d, % 3d bytes: % x", nonce.ctr, n, chunk[:n]))
			chunk = Pad(chunk[:n], size, last)
			//ct := c.aead.Seal(nil, nonce.Next(), chunk[:n], nil)
			nw, err := w.Write(chunk)
			nOut += int64(nw)
			if err != nil {
				errOut = err
				return
			}
		}
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		} else if err != nil {
			errOut = err
			return
		}

	}
	return

}

func (c *crypter) Decrypt(w io.Writer, r io.Reader) (nOut int64, errOut error) {

	// buffered reader
	br := bufio.NewReader(r)

	// encrypted chunk size and counter
	size := CHUNKSIZE // + c.aead.Overhead()
	chunk := make([]byte, size)
	nonce := NewNonceCounter()

	for {

		n, err := io.ReadFull(br, chunk[:size])
		if n > 0 {
			stderr(sfmt("input  chunk % 3d, % 3d bytes: % x", nonce.ctr, n, chunk[:n]))
			unp, last := Unpad(chunk[:n])
			nw, err := w.Write(unp)
			nOut += int64(nw)
			if err != nil {
				errOut = err
				return
			}
			if last {
				break
			}
		}
		if err != nil {
			errOut = err
			return
		}

	}
	return

}
