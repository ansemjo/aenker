package aenker

import (
	"io"

	ce "github.com/ansemjo/aenker/error"
	padding "github.com/ansemjo/aenker/padding"
)

// ErrExtraData indicates that there is extra data appended to the ciphertext
const ErrExtraData = ce.ConstError("aenker: extraneous data after ciphertext")

// Encrypt reads from the given Reader, chunks the plaintext into
// equal parts, pads and encrypts them with an AEAD and writes ciphertext
// to the given Writer
func (a *Aenker) Encrypt(w io.Writer, r io.Reader) (lengthWritten int64, error error) {

	size, buf, chunk, nonce := a.initializeMode(r, encrypt) // initialize needed structures
	for {

		nr, rErr := io.ReadFull(buf, chunk[:size-1])                // read and leave room for pad
		final := eof(buf)                                           // check if this is the last chunk
		if nr > 0 && (rErr == nil || rErr == io.ErrUnexpectedEOF) { // if there is data and no unusual error

			padding.AddPadding(chunk[:nr], final)            // add padding to plaintext
			ct := a.aead.Seal(nil, nonce.Next(), chunk, nil) // encrypt padded data, increment nonce
			nw, wErr := w.Write(ct)                          // write ciphertext to writer
			lengthWritten += int64(nw)                       // update output length
			if wErr != nil {                                 // an error occurred during write
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

	_, buf, chunk, nonce := a.initializeMode(r, decrypt) // initialize needed structures
	for {

		nr, rErr := io.ReadFull(buf, chunk)                         // read what is supposed to be a ciphertext chunk
		more := !eof(buf)                                           // check if there is more data
		if nr > 0 && (rErr == nil || rErr == io.ErrUnexpectedEOF) { // if there is data and no unusual error

			pt, cErr := a.aead.Open(nil, nonce.Next(), chunk[:nr], nil) // decrypt and verify data, increment nonce
			if cErr != nil {                                            // an error occurred during decryption
				error = cErr
				return
			}

			dlen, final := padding.RemovePadding(pt) // get data length after padding removal
			nw, wErr := w.Write(pt[:dlen])           // write plaintext slice to writer
			lengthWritten += int64(nw)               // update output length
			if wErr != nil {                         // an error occurred during write
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
