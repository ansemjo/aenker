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

// Header is the struct that is serialized at the beginning of encrypted
// aenker files. It is required to derive the shared key upon decryption and
// its serialization is used as associated data during chunk sealing/opening.
//
// Rather than writing or opening the header manually, use NewWriter or NewReader
// to start encrypting a new file or decrypt a previously encrypted file.
type Header struct {
	Magic     [8]byte
	Salt      [8]byte
	Ephemeral [32]byte
}

// Magic is the magic bytes string that is used to identify aenker files.
// The two bytes after 'aenker' are the first two bytes of its Blake2b hash:
//  >>> hashlib.blake2b(b'aenker').digest()[:2]
//  b'\xe7\x9e'
const Magic = "aenker\xe7\x9e"

// Keyinfo is used as context info for HKDF during key derivation.
const Keyinfo = "aenker elliptic"

func writeNewHeader(writer io.Writer, peer *[32]byte) (key, head []byte, err error) {

	// create new header struct and copy magic bytes
	header := &Header{}
	copy(header.Magic[:], []byte(Magic))

	// new random salt
	_, err = io.ReadFull(rand.Reader, header.Salt[:])
	if err != nil {
		return
	}

	// new ephemeral secret key
	_, err = io.ReadFull(rand.Reader, header.Ephemeral[:])
	if err != nil {
		return
	}
	// derive shared key for chunkstream
	key = keyderivation.Elliptic(&header.Ephemeral, peer, header.Salt[:], Keyinfo)

	// replace ephemeral with its public key
	header.Ephemeral = *keyderivation.Public(&header.Ephemeral)

	// create a small buffer to hold the written header
	buf := bytes.NewBuffer(make([]byte, 0, 48))

	// write header to writer and a buffer
	tee := io.MultiWriter(buf, writer)
	err = binary.Write(tee, binary.BigEndian, header)
	if err != nil {
		key = nil
		return
	}

	// return derived key and written header
	return key, buf.Bytes(), err

}

func openHeader(reader io.Reader, private *[32]byte) (key, head []byte, err error) {

	// create a small buffer to hold the read header
	buf := bytes.NewBuffer(make([]byte, 0, 48))

	// create new header struct
	header := &Header{}

	// read, buffer and decode header
	tee := io.TeeReader(reader, buf)
	err = binary.Read(tee, binary.BigEndian, header)
	if err != nil {
		return
	}

	// check magic bytes, public data so no constant-time implementation
	if bytes.Compare(header.Magic[:], []byte(Magic)) != 0 {
		err = errors.New("unknown magic bytes")
	}

	// derive shared key for chunkstream
	key = keyderivation.Elliptic(private, &header.Ephemeral, header.Salt[:], Keyinfo)

	// return derived key and read header
	return key, buf.Bytes(), err

}
