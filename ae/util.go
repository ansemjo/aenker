// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package ae

import (
	"bufio"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// find out if the reader is exhausted by
// peeking ahead one byte
func eof(r *bufio.Reader) bool {
	_, err := r.Peek(1)
	return err == io.EOF
}

// any non-nil error is a fatal failure.
// print error to stderr and exit
func fatal(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// return bytes from system randomness
func randomBytes(size int) (bytes []byte) {
	bytes = make([]byte, size)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(fmt.Sprintf("could not read %d random bytes", size))
	}
	return
}

// print bytes in slice to stderr
func debugBytes(label string, slice []byte) {
	fmt.Fprintf(os.Stderr, "%s: % x\n", label, slice)
}

// simple binary encoding of uint32's
func itob(i uint32) (b []byte) {
	b = make([]byte, 4)
	binary.LittleEndian.PutUint32(b, i)
	return
}
func btoi(b []byte) (i uint32) {
	return binary.LittleEndian.Uint32(b)
}
