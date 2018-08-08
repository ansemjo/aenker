package main

import (
	"crypto/cipher"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/crypto/chacha20poly1305"
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

func (c *crypter) Encrypt(w io.Writer, r io.Reader) (nOut int64, errOut error) {

	// internal chunk size and counter
	chunk := make([]byte, 32)
	nonce := NewNonceCounter()
	fmt.Println(len(nonce.Next()))

	for {

		n, err := io.ReadFull(r, chunk)
		if n > 0 {
			ct := c.aead.Seal(nil, nonce.Next(), chunk[:n], nil)
			nw, err := w.Write(ct)
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

func main() {

	reader := strings.NewReader("Clear is better than clever")
	writer := os.Stdout
	fmt.Fprintln(os.Stderr, "reader has", reader.Len(), "bytes")

	n, err := newCrypter().Encrypt(writer, reader)
	fmt.Fprintln(os.Stderr, "wrote", n, "bytes")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
