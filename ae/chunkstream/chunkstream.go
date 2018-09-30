// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package chunkstream

import (
	"bufio"
	"crypto/cipher"
	"errors"
	"io"

	"github.com/ansemjo/aenker/ae/padding"
)

// ErrExtraData indicates that there is extra data appended to the ciphertext
var ErrExtraData = errors.New("aenker: extraneous data after ciphertext")

// associated data for the chunk encryption
const chunkEncryptionAD = "Aenker Chunk"

// wether we are encrypting or decrypting
type mode int

const encrypt mode = 0
const decrypt mode = 1

type chunkStreamer struct {
	size   int
	aead   cipher.AEAD
	reader *bufio.Reader
	ctr    *nonceCounter
	ad     []byte
}

// perform some common initialization for encryption and decryption
func (ae *Aenker) newChunkStreamer(r io.Reader, mode mode) (
	streamer *chunkStreamer, chunkbuf []byte, err error) {

	s := &chunkStreamer{}         // new chunk streamer struct
	s.reader = bufio.NewReader(r) // add buffered reader
	s.ctr = newNonceCounter()     // add nonce counter
	s.ad = append([]byte(chunkEncryptionAD), itob(uint32(ae.chunksize))...)

	s.aead, err = ae.cipher.New(*ae.mek) // init AEAD with media encryption key
	if err != nil {
		return
	}

	if mode == encrypt { // depending on mode, calculate chunk size
		s.size = ae.chunksize
	} else {
		s.size = ae.chunksize + s.aead.Overhead()
	}

	streamer = s                    // assign created streamer pointer
	chunkbuf = make([]byte, s.size) // allocate slice with correct capacity

	return

}

// EncryptChunkStream reads from the given Reader, chunks the plaintext into
// equal parts, pads and encrypts them with an AEAD and writes ciphertext
// to the given Writer
func (ae *Aenker) encryptChunkStream(w io.Writer, r io.Reader) (lengthWritten uint64, error error) {

	s, chunk, error := ae.newChunkStreamer(r, encrypt) // initialize streamer structure and allocate memory
	if error != nil {
		return
	}

	for {

		nr, rErr := io.ReadFull(s.reader, chunk[:s.size-1])         // read next chunk and leave room for pad
		final := eof(s.reader)                                      // check if this is the last chunk
		if nr > 0 && (rErr == nil || rErr == io.ErrUnexpectedEOF) { // if there is data and no unusual error

			chunk = chunk[:nr]                                // truncate to read data
			padding.AddPadding(&chunk, final, ae.chunksize)   // add padding to plaintext
			ct := s.aead.Seal(nil, s.ctr.Next(), chunk, s.ad) // encrypt padded data, increment nonce
			nw, wErr := w.Write(ct)                           // write ciphertext to writer
			lengthWritten += uint64(nw)                       // update output length
			if wErr != nil {                                  // an error occurred during write
				error = wErr
				return
			}

		} else { // possibly fatal error occurred during read
			error = rErr
			return
		}
		if final { // this was the last chunk, we're done
			break
		}
	}
	return

}

// DecryptChunkStream read ciphertext from the given reader, chunks the AEAD-encrypted blocks,
// attempts to decrypt and verify them and finally writes the original plaintext
// to the given Writer
func (ae *Aenker) decryptChunkStream(w io.Writer, r io.Reader) (lengthWritten uint64, error error) {

	s, chunk, error := ae.newChunkStreamer(r, decrypt) // initialize streamer structure and allocate memory
	if error != nil {
		return
	}

	for {

		nr, rErr := io.ReadFull(s.reader, chunk)                    // read what is supposed to be a ciphertext chunk
		more := !eof(s.reader)                                      // check if there is more data
		if nr > 0 && (rErr == nil || rErr == io.ErrUnexpectedEOF) { // if there is data and no unusual error

			pt, cErr := s.aead.Open(nil, s.ctr.Next(), chunk[:nr], s.ad) // decrypt and verify data, increment nonce
			if cErr != nil {                                             // an error occurred during decryption
				error = cErr
				return
			}

			final := padding.RemovePadding(&pt) // get data length after padding removal
			nw, wErr := w.Write(pt)             // write plaintext slice to writer
			lengthWritten += uint64(nw)         // update output length
			if wErr != nil {                    // an error occurred during write
				error = wErr
				return
			}
			if final { // this was the last chunk, we're done
				if more { // but the was extraneous data
					error = ErrExtraData // set informative error
				}
				break
			}
		} else { // possibly fatal error occurred during read
			error = rErr
			return
		}

	}
	return

}
