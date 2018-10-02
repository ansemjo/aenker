// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package chunkstream

import (
	"bytes"
	"io"

	"github.com/ansemjo/aenker/padding"
)

type chunkReader struct {
	chipherer *chunkCipherer
	buf       *bytes.Buffer
	chunksize int
	reader    io.Reader
	err       error
}

// NewReader instantiates a new authenticated cipher from NewAEAD with the given key and
// returns a Reader. Any reads from that will read and buffer an appropriate amount of encrypted
// data to return the next chunk before being decrypted and authenticated. Only successfully
// authenticated data is ever returned.
//
// Do not increase the chunksize manually to compensate for AEAD overhead, the chunkCipherer within
// will do that automatically. I.e. if you encrypted with chunksize=2048 you need to decrypt with
// chunksize=2048.
func NewReader(r io.Reader, key, info []byte, chunksize int) (io.Reader, error) {

	cr := &chunkReader{reader: r}
	var err error

	cr.chipherer, err = newChunkCipherer(key, info)

	cr.chunksize = chunksize + cr.chipherer.cipher.Overhead()

	if err == nil {
		cr.buf = bytes.NewBuffer(make([]byte, 0, chunksize))
	}

	return cr, err

}

func (cr *chunkReader) Read(p []byte) (n int, err error) {

	// previous errors
	if cr.err != nil {
		return 0, cr.err
	}
	// save error for future calls upon exit
	defer func() {
		if err != nil {
			cr.err = err
		}
	}()

	// decrypt more data
	if cr.buf.Len() == 0 {
		err = cr.open()
		if err == io.EOF {
			cr.err = err
		} else if err != nil {
			return
		}
	}

	return cr.buf.Read(p)

}

func (cr *chunkReader) open() (err error) {

	// TODO: direct copy to second internal buffer with io.CopyN ?
	chunk := make([]byte, cr.chunksize)
	_, err = io.ReadFull(cr.reader, chunk)
	if err != nil {
		return
	}

	chunk, err = cr.chipherer.Open(chunk)
	if err != nil {
		return
	}
	final := padding.Remove(&chunk)
	_, err = cr.buf.Write(chunk)
	if err != nil {
		return
	}
	if final {
		err = io.EOF
	}
	return

}
