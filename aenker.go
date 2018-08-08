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
	reader io.Reader
	aead   cipher.AEAD
	nonce  []byte
}

func newCrypter(reader io.Reader) *crypter {

	// init zero key and nonce
	key := make([]byte, chacha20poly1305.KeySize)
	nonce := make([]byte, chacha20poly1305.NonceSize)

	aead, err := chacha20poly1305.New(key)
	if err != nil {
		panic("AEAD init failed")
	}

	return &crypter{
		reader: reader,
		aead:   aead,
		nonce:  nonce,
	}
}

func (c *crypter) WriteTo(w io.Writer) (nOut int64, errOut error) {

	buf := make([]byte, 32)
	for {

		n, err := io.ReadFull(c.reader, buf)
		if n > 0 {
			ct := c.aead.Seal(nil, c.nonce, buf[:n], nil)
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

	n, err := newCrypter(reader).WriteTo(writer)
	fmt.Fprintln(os.Stderr, "wrote", n, "bytes")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
