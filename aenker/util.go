package aenker

import (
	"bufio"
	"io"
)

// find out if the reader is exhausted by
// peeking ahead one byte
func eof(r *bufio.Reader) bool {
	_, err := r.Peek(1)
	return err == io.EOF
}

// perform some common initialization for encryption and decryption
func (a *Aenker) initializeMode(r io.Reader, mode mode) (
	size int, bufferedReader *bufio.Reader, chunk []byte, nonce *Nonce) {

	if mode == encrypt {
		size = a.chunksize
	} else {
		size = a.chunksize + a.aead.Overhead()
	}
	//? TODO: does NewReaderSize make sense? apply size constraints?
	bufferedReader = bufio.NewReader(r)
	chunk = make([]byte, size)
	nonce = NewNonce()
	return

}
