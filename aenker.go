package main

import (
	"os"
	"strings"
)

func main() {

	reader := strings.NewReader("Clear is better than clever")
	writer := os.Stdout

	stderr(sfmt("reader has %d bytes", reader.Len()))
	n, err := newCrypter().Encrypt(writer, reader)
	stderr(sfmt("wrote %d bytes", n))
	if err != nil {
		stderr(err)
	}
}
