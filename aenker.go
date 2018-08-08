package main

import (
	"bytes"
	"os"
	"strings"
)

func main() {

	reader := strings.NewReader("Clear is better than clever.\n")
	// stderr(sfmt("reader has %d bytes", reader.Len()))

	writer := os.Stdout
	var buffer bytes.Buffer

	_, err := newCrypter().Encrypt(&buffer, reader)
	// stderr(sfmt("wrote %d bytes", n))
	if err != nil {
		stderr(err)
	}

	//io.Copy(os.Stderr, &buffer)
	newCrypter().Decrypt(writer, &buffer)

}
