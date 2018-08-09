package aenker

import (
	"bufio"
	"crypto/cipher"
	"io"

	constError "github.com/ansemjo/aenker/error"
	padding "github.com/ansemjo/aenker/padding"
	"golang.org/x/crypto/chacha20poly1305"
)

// wether we are encrypting or decrypting
type mode int

const encrypt mode = 0
const decrypt mode = 1

const (

	// ChunkSize is the amount of plaintext that is encrypted per chunk
	ChunkSize = 256

	// ErrExtraData indicates that there is extra data appended to the ciphertext
	ErrExtraData = constError.Error("aenker: extraneous data after ciphertext")
)

// Aenker is a struct that supports a sort of "streamed" AEAD usage,
// where the plaintext is split into equal parts and encrypted with
// an interMAClib-like construction
type Aenker struct {
	aead cipher.AEAD
}

// NewAenker return a new Aenker with initialized cipher
func NewAenker(key []byte) (ae *Aenker, keyerror error) {
	aead, err := chacha20poly1305.New(key)
	return &Aenker{aead}, err
}

// find out if the reader is exhausted by
// peeking ahead one byte
func eof(r *bufio.Reader) bool {
	_, err := r.Peek(1)
	return err == io.EOF
}

// perform some common initialization for encryption and decryption
func (a *Aenker) initCommon(r io.Reader, mode mode) (
	size int, bufferedReader *bufio.Reader, chunk []byte, nonce *NonceCounter) {

	if mode == encrypt {
		size = ChunkSize
	} else {
		size = ChunkSize + a.aead.Overhead()
	}
	//? TODO: does NewReaderSize make sense? apply size constraints?
	bufferedReader = bufio.NewReader(r)
	chunk = make([]byte, size)
	nonce = NewNonceCounter()
	return

}

// Encrypt reads from the given Reader, chunks the plaintext into
// equal parts, pads and encrypts them with an AEAD and writes ciphertext
// to the given Writer
func (a *Aenker) Encrypt(w io.Writer, r io.Reader) (lengthWritten int64, error error) {

	size, buf, chunk, nonce := a.initCommon(r, encrypt) // initialize needed structures
	for {

		nr, rErr := io.ReadFull(buf, chunk[:size-1])                // read and leave room for pad
		final := eof(buf)                                           // check if this is the last chunk
		if nr > 0 && (rErr == nil || rErr == io.ErrUnexpectedEOF) { // if there is data and no unusual error

			// TODO: work on original slice, so chunk[:size] is not needed
			chunk = padding.Pad(chunk[:nr], size, final)            // add padding to plaintext
			ct := a.aead.Seal(nil, nonce.Next(), chunk[:size], nil) // encrypt padded data, increment nonce
			nw, wErr := w.Write(ct)                                 // write ciphertext to writer
			lengthWritten += int64(nw)                              // update output length
			if wErr != nil {                                        // an error occurred during write
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

// Decrypt read ciphertext from the given reader, chunks the AEAD-encrypted blocks,
// attempts to decrypt and verify them and finally writes the original plaintext
// to the given Writer
func (a *Aenker) Decrypt(w io.Writer, r io.Reader) (lengthWritten int64, error error) {

	_, buf, chunk, nonce := a.initCommon(r, decrypt) // initialize needed structures
	for {

		nr, rErr := io.ReadFull(buf, chunk)                         // read what is supposed to be a ciphertext chunk
		more := !eof(buf)                                           // check if there is more data
		if nr > 0 && (rErr == nil || rErr == io.ErrUnexpectedEOF) { // if there is data and no unusual error

			pt, cErr := a.aead.Open(nil, nonce.Next(), chunk[:nr], nil) // decrypt and verify data, increment nonce
			if cErr != nil {                                            // an error occurred during decryption
				error = cErr
				return
			}

			pt, final := padding.Unpad(pt) // remove padding from plaintext, pad indicates if this is the last chunk
			nw, wErr := w.Write(pt)        // write plaintext to writer
			lengthWritten += int64(nw)     // update output length
			if wErr != nil {               // an error occurred during write
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
