package chunkstream

import (
	"bytes"
	"io"

	"github.com/ansemjo/aenker/ae/padding"
)

type ChunkReader struct {
	chipherer *chunkCipherer
	buf       *bytes.Buffer
	chunksize int
	reader    io.Reader
	err       error
}

func NewChunkReader(r io.Reader, key, info []byte, chunksize int) (io.Reader, error) {

	cr := &ChunkReader{reader: r}
	var err error

	cr.chipherer, err = newChunkCipherer(key, info)

	cr.chunksize = chunksize + cr.chipherer.cipher.Overhead()

	if err == nil {
		cr.buf = bytes.NewBuffer(make([]byte, 0, chunksize))
	}

	return cr, err

}

func (cr *ChunkReader) Read(p []byte) (n int, err error) {

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
	//return io.ReadFull(cr.reader, p)

}

func (cr *ChunkReader) open() (err error) {

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
	final := padding.RemovePadding(&chunk)
	_, err = cr.buf.Write(chunk)
	if err != nil {
		return
	}
	if final {
		err = io.EOF
	}
	return

}
