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

// TODO: create WriteCloser and Reader interfaces

// Options is a configuration object required to instantiate a chunkStreamer
type Options struct {
	AEAD      func([]byte) (cipher.AEAD, error)
	Key       []byte
	Info      []byte
	Reader    io.Reader
	Writer    io.Writer
	Chunksize int
}

// chunkStream is the internal state of a chunkStreamer
// TODO: do we need two seperate option structs?
type chunkStream struct {
	aead   cipher.AEAD
	info   []byte
	size   int
	reader *bufio.Reader
	writer io.Writer
	ctr    *nonceCounter
}

// perform some common initialization for encryption and decryption
// TODO: properly error out if any required field is nil
func newChunkStream(opts Options, encrypt bool) (stream *chunkStream, chunk []byte, err error) {

	s := &chunkStream{
		reader: bufio.NewReader(opts.Reader), // buffered reader
		writer: opts.Writer,                  // copy writer
		info:   opts.Info,                    // additional data for aead
	}

	s.aead, err = opts.AEAD(opts.Key) // init aead with media encryption key
	if err != nil {
		return
	}

	s.ctr = newNonceCounter(s.aead.NonceSize()) // nonce counter

	if encrypt { // depending on mode, calculate chunk size
		s.size = opts.Chunksize
	} else {
		s.size = opts.Chunksize + s.aead.Overhead()
	}

	stream = s                   // assign created streamer pointer
	chunk = make([]byte, s.size) // allocate slice with correct capacity

	return

}

// Encrypt reads from the given Reader, chunks the plaintext into
// equal parts, pads and encrypts them with an AEAD and writes ciphertext
// to the given Writer
func Encrypt(opts Options) (lengthWritten uint64, err error) {

	s, chunk, err := newChunkStream(opts, true)
	if err != nil {
		return
	}

	for {

		nr, rErr := io.ReadFull(s.reader, chunk[:s.size-1])         // read next chunk and leave room for pad
		final := eof(s.reader)                                      // check if this is the last chunk
		if nr > 0 && (rErr == nil || rErr == io.ErrUnexpectedEOF) { // if there is data and no unusual error

			chunk = chunk[:nr]                                  // truncate to read data
			padding.AddPadding(&chunk, final, s.size)           // add padding to plaintext
			ct := s.aead.Seal(nil, s.ctr.Next(), chunk, s.info) // encrypt padded data, increment nonce
			nw, wErr := s.writer.Write(ct)                      // write ciphertext to writer
			lengthWritten += uint64(nw)                         // update output length
			if wErr != nil {                                    // an error occurred during write
				err = wErr
				return
			}

		} else { // possibly fatal error occurred during read
			err = rErr
			return
		}
		if final { // this was the last chunk, we're done
			break
		}
	}
	return

}

// Decrypt read ciphertext from the given reader, chunks the AEAD-encrypted blocks,
// attempts to decrypt and verify them and finally writes the original plaintext
// to the given Writer
func Decrypt(opts Options) (lengthWritten uint64, err error) {

	s, chunk, err := newChunkStream(opts, false)
	if err != nil {
		return
	}

	for {

		nr, rErr := io.ReadFull(s.reader, chunk)                    // read what is supposed to be a ciphertext chunk
		more := !eof(s.reader)                                      // check if there is more data
		if nr > 0 && (rErr == nil || rErr == io.ErrUnexpectedEOF) { // if there is data and no unusual error

			pt, cErr := s.aead.Open(nil, s.ctr.Next(), chunk[:nr], s.info) // decrypt and verify data, increment nonce
			if cErr != nil {                                               // an error occurred during decryption
				err = cErr
				return
			}

			final := padding.RemovePadding(&pt) // get data length after padding removal
			nw, wErr := s.writer.Write(pt)      // write plaintext slice to writer
			lengthWritten += uint64(nw)         // update output length
			if wErr != nil {                    // an error occurred during write
				err = wErr
				return
			}
			if final { // this was the last chunk, we're done
				if more { // but the was extraneous data
					err = ErrExtraData // set informative error
				}
				break
			}
		} else { // possibly fatal error occurred during read
			err = rErr
			return
		}

	}
	return

}

// ErrExtraData is an informative error to signal that there was extraneous data after last chunk
var ErrExtraData = errors.New("aenker: extraneous data after ciphertext")
