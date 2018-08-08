package main

import (
	"bytes"
	"io"
	"os"
	"strings"
)

func main() {

	reader := strings.NewReader("Clear is better than clever.\n")
	// stderr(sfmt("reader has %d bytes", reader.Len()))

	writer := os.Stdout
	var midbuf bytes.Buffer
	var outbuf bytes.Buffer

	_, err := newCrypter().Encrypt(&midbuf, reader)
	// stderr(sfmt("wrote %d bytes", n))
	if err != nil {
		panic(err)
	}

	//io.Copy(os.Stderr, &buffer)
	_, err = newCrypter().Decrypt(&outbuf, &midbuf)
	if err != nil {
		panic(err)
	}

	io.Copy(writer, &outbuf)

}
