package chunkstream

import (
	"bytes"
	"io"

	"github.com/ansemjo/aenker/padding"
)

type chunkWriter struct {
	chipherer *chunkCipherer
	buf       *bytes.Buffer
	chunksize int
	writer    io.Writer
	err       error
}

// NewWriter instantiates a new authenticated cipher from NewAEAD with the given key and
// returns a WriteCloser. Any writes to that will be split into small chunks and is then
// encrypted and authenticated individually before being written to the passed Writer.
//
// You MUST call Close() when you're done to ensure the final chunk is written.
//
// You MUST use a unique key because internally a simple incrementing counter is used as
// a nonce, so two streams encrypted with the same key will compromise confidentiality!
func NewWriter(w io.Writer, key, info []byte, chunksize int) (io.WriteCloser, error) {

	cw := &chunkWriter{chunksize: chunksize, writer: w}
	var err error

	cw.chipherer, err = newChunkCipherer(key, info)

	if err == nil {
		cw.buf = bytes.NewBuffer(make([]byte, 0, chunksize))
	}

	return cw, err

}

func (cw *chunkWriter) Write(data []byte) (n int, err error) {

	// previous errors
	if cw.err != nil {
		return 0, cw.err
	}
	// save error for future calls upon exit
	defer func() {
		if err != nil {
			cw.err = err
		}
	}()

	// while there is data
	for len(data) > 0 {

		// maximum needed for next chunk, capped to available data
		need := min(cw.chunksize-cw.buf.Len()-1, len(data))

		// write more data to buffer
		nb, err := cw.buf.Write(data[:need])
		n += nb
		if err != nil {
			return n, err
		}
		// reslice input data
		data = data[nb:]

		// process chunk if there is enough data
		if cw.buf.Len() >= cw.chunksize-1 {
			err = cw.seal(false)
			if err != nil {
				return n, err
			}
		}

	}

	return

}

func (cw *chunkWriter) seal(final bool) (err error) {

	chunk := cw.buf.Next(cw.chunksize - 1)
	padding.AddPadding(&chunk, final, cw.chunksize) // add padding to plaintext
	ct := cw.chipherer.Seal(chunk)                  // encrypt padded data, increment nonce
	_, err = cw.writer.Write(ct)                    // write ciphertext to writer
	return

}

func (cw *chunkWriter) Close() (err error) {
	return cw.seal(true)
}

// return smaller int
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
