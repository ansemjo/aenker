// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package ae

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"io"

	"github.com/ansemjo/aenker/keyderivation"
)

// the two bytes after 'aenker' are the first two bytes of its blake2b hash:
// >>> hashlib.blake2b(b'aenker').digest()[:2]
// b'\xe7\x9e'
const filemagic = "aenker\xe7\x9e"

// used as context info for hkdf during key derivation
const keyinfo = "aenker elliptic"

func WriteNewHeader(writer io.Writer, peer *[32]byte) (key []byte, err error) {

	// create new header struct and copy magic bytes
	header := &header{}
	copy(header.magic[:], []byte(filemagic))

	// new random salt
	_, err = io.ReadFull(rand.Reader, header.salt[:])
	if err != nil {
		return
	}

	// new ephemeral secret key
	_, err = io.ReadFull(rand.Reader, header.ephemeral[:])
	if err != nil {
		return
	}
	// derive shared key for chunkstream
	key = keyderivation.Elliptic(&header.ephemeral, peer, header.salt[:], keyinfo)

	// replace ephemeral with its public key
	header.ephemeral = *keyderivation.Public(&header.ephemeral)

	// write header
	err = binary.Write(writer, binary.BigEndian, header)
	if err != nil {
		key = nil
		return
	}

	// should have written 48 bytes if successful
	return key, err

}

func OpenHeader(reader io.Reader, private *[32]byte) (key []byte, err error) {

	// create new header struct
	header := &header{}

	// read and decode header
	err = binary.Read(reader, binary.BigEndian, header)
	if err != nil {
		return
	}

	// check magic bytes, public data so no constant-time implementation
	if bytes.Compare(header.magic[:], []byte(filemagic)) != 0 {
		err = errors.New("unknown magic bytes")
	}

	// derive shared key for chunkstream
	key = keyderivation.Elliptic(private, &header.ephemeral, header.salt[:], keyinfo)

	return

}
