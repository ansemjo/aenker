// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

// Package ae provides functions to encrypt data written to a Writer or decrypt data read from
// a Reader. It implements a chunked ECIES [0] with Curve25519 [1] and ChaCha20Poly1305 [2] internally.
// Chunking and individual padding is done similarly to InterMAC [3].
//
//  [0]: https://en.wikipedia.org/wiki/Integrated_Encryption_Scheme
//  [1]: https://cr.yp.to/ecdh.html / https://tools.ietf.org/html/rfc7748#section-4.1
//  [2]: https://tools.ietf.org/html/rfc7539
//  [3]: https://rwc.iacr.org/2018/Slides/Hansen.pdf
package ae

import (
	"io"

	"github.com/ansemjo/aenker/chunkstream"
)

// TODO: add links to diagrams when they are finalised and added to the repository.

// Chunksize is the artificial chunksize used for the ChunkStream to write small
// encrypted files of exactly 2 kB. This is also a good value to reduce the losses
// through padding and overhead to < 1% on files larger than 1 MB.
const Chunksize = 1984 // big brother is watching you

// NewWriter derives an ephemeral shared key with the given Curve25519 public key,
// writes a header to the provided Writer and then returns a ChunkWriter, which will encrypt
// any written data.
//
// Don't forget to call .Close() when you're done, otherwise the final chunk will never be
// written. This does NOT close the writer that was originally passed though, i.e. if you passed
// a file, you need to close that seperately!
func NewWriter(w io.Writer, public *[32]byte) (cw io.WriteCloser, err error) {

	// write new header and derive key
	key, head, err := writeNewHeader(w, public)
	if err != nil {
		return
	}

	return chunkstream.NewWriter(w, key, head, Chunksize)

}

// NewReader tries to open the given Reader, decode the header and derive a shared key
// with your private and the decoded ephemeral public key. It then returns a ChunkReader, which
// will transparently decrypt data upon calling .Read().
//
// Please note that opening a valid header will succeed even if you provide the wrong private
// key as the header itself is not MAC'ed. It is however used as associated data in the chunks, so
// decryption will fail upon the first call to .Read().
func NewReader(r io.Reader, private *[32]byte) (cr io.Reader, err error) {

	// open header and derive key
	key, head, err := openHeader(r, private)
	if err != nil {
		return
	}

	return chunkstream.NewReader(r, key, head, Chunksize)

}
