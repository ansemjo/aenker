// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package ae

import (
	"fmt"
	"io"
)

// "aenker" + version + kdf mode + chunksize + nonce
const aenkerV2HeaderLen = 6 + 1 + 1 + 2 + 32

func (ae *Aenker2) OpenHeader(r io.Reader) (err error) {
	this := "ae.openHeader"
	E := func(s string) error { return errfmt(this, s) }

	// read entire header at once
	header, err := readBytes(r, aenkerV2HeaderLen)
	if err != nil {
		return fmt.Errorf("%s: short header read: %s", this, err)
	}

	// check magic
	m := header[:6]
	if magic != string(m) {
		return E("wrong magic bytes")
	}

	// ignore version for now

	// parse key derivation mode
	k := header[7]
	if k > kdfModeMax {
		return E("unknown key derivation mode")
	}
	ae.kdf = kdfMode(k)

	// read chunksize
	cs := header[8:10]
	ae.chunksize = btou16(cs)

	// read nonce
	n := header[10:42]
	ae.nonce = n

	return

}
