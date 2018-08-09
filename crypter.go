package main

import (
	"bufio"
	"crypto/cipher"
	"io"

	"golang.org/x/crypto/chacha20poly1305"
)

// wether we are encrypting or decrypting
type mode int

const encrypt mode = 0
const decrypt mode = 1

const (
	// CHUNKSIZE is the amount of plaintext that is encrypted per chunk
	CHUNKSIZE = 256
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

// find out if the reader is exhausted by
// peeking ahead one byte
func isLast(r *bufio.Reader) bool {
	_, err := r.Peek(1)
	return err == io.EOF
}

// Encrypt reads from the given Reader, chunks the plaintext into
// equal parts, pads and encrypts them with an AEAD and writes ciphertext
// to the given Writer
func (a *Aenker) Encrypt(w io.Writer, r io.Reader) (nOut int64, error error) {

	size, buf, chunk, nonce := a.initCommon(r, encrypt)
	for {

		n, err := io.ReadFull(buf, chunk[:size-1])
		last := isLast(buf)
		if n > 0 {
			//stderr(sfmt("output chunk % 3d, % 3d bytes: % x", nonce.ctr, n, chunk[:n]))
			chunk = Pad(chunk[:n], size, last)
			ct := a.aead.Seal(nil, nonce.Next(), chunk[:size], nil)
			//stderr(sfmt("cipher chunk % 3d, % 3d bytes: % x", nonce.ctr-1, len(ct), ct))
			nw, err := w.Write(ct)
			nOut += int64(nw)
			if err != nil {
				error = err
				return
			}
		}
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		} else if err != nil {
			error = err
			return
		}

	}
	return

}

// perform some common initialization for encryption and decryption
func (a *Aenker) initCommon(r io.Reader, mode mode) (
	size int, bufferedReader *bufio.Reader, chunk []byte, nonce *NonceCounter) {

	if mode == encrypt {
		size = CHUNKSIZE
	} else {
		size = CHUNKSIZE + a.aead.Overhead()
	}
	//? TODO: does NewReaderSize make sense? apply size constraints?
	bufferedReader = bufio.NewReader(r)
	chunk = make([]byte, size)
	nonce = NewNonceCounter()
	return

}

func (a *Aenker) Decrypt(w io.Writer, r io.Reader) (nOut int64, errOut error) {

	_, buf, chunk, nonce := a.initCommon(r, decrypt)
	for {

		n, err := io.ReadFull(buf, chunk)
		if n > 0 {
			//stderr(sfmt("input  chunk % 3d, % 3d bytes: % x", nonce.ctr, n, chunk[:n]))
			pt, err := a.aead.Open(nil, nonce.Next(), chunk[:n], nil)
			if err != nil {
				errOut = err
				return
			}
			//stderr(sfmt("plain  chunk % 3d, % 3d bytes: % x", nonce.ctr-1, len(pt), pt))
			unp, last := Unpad(pt)
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
