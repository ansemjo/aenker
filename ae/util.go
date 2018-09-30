// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package ae

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

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

// binary encoding of uint16
func u16tob(u uint16) (b []byte) {
	b = make([]byte, 2)
	binary.LittleEndian.PutUint16(b, u)
	return
}
func btou16(b []byte) (u uint16) {
	return binary.LittleEndian.Uint16(b)
}

// read specified number of bytes from reader
func readBytes(r io.Reader, n int) (b []byte, err error) {
	b = make([]byte, n)
	_, err = io.ReadFull(r, b)
	return
}

// read a single byte from reader
func readByte(r io.Reader) (b byte, err error) {
	by, err := readBytes(r, 1)
	if err != nil {
		return
	}
	return by[0], nil
}

// concatenate string and error as new error
func errfmt(str, err string) error {
	return fmt.Errorf("%s: %s", str, err)
}
