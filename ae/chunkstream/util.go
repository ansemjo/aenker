package chunkstream

import (
	"bufio"
	"io"
)

// find out if the reader is exhausted by
// peeking ahead one byte
func eof(r *bufio.Reader) bool {
	_, err := r.Peek(1)
	return err == io.EOF
}

// return smaller int
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
