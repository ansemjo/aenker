// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package ae

import (
	"io"
	"os"
)

// Read a file from disk and output the content to stdout.
func ExampleNewReader_file() {

	// open file on disk
	file, err := os.Open("secrets.ae")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// dummy private key
	key := new([32]byte)

	// open aenker reader
	ae, err := NewReader(file, key)
	if err != nil {
		panic(err)
	}

	// decrypt and copy to stdout
	io.Copy(os.Stdout, ae)

}

// Encrypt data passed in os.Stdin and write to os.Stdout (like a unix pipe).
func ExampleNewWriter_pipe() {

	// dummy public key
	peer := new([32]byte)

	// open aenker writer
	ae, err := NewWriter(os.Stdout, peer)
	if err != nil {
		panic(err)
	}
	// close writer when you're done to ensure last chunk is written!
	defer ae.Close()

	// copy input to output while transparently encrypting
	io.Copy(ae, os.Stdin)

}
