// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package chunkstream

import "encoding/binary"

// nonceCounter is a 64 bit counter which is monotonically incremented
// and outputs a slice on .Next() for AEAD usage
type nonceCounter struct {
	ctr   uint64
	nonce []byte
	size  int
}

// newNonceCounter returns a new counter starting at 0 that will output
// slices with length 'size' upon calling Next().
func newNonceCounter(size int) *nonceCounter {
	return &nonceCounter{
		nonce: make([]byte, 32),
		size:  size,
	}
}

// Next outputs the current counter value as a slice and then increments
// the internal counter.
func (nc *nonceCounter) Next() (nonce []byte) {
	binary.LittleEndian.PutUint64(nc.nonce, nc.ctr)
	nc.ctr++
	return nc.nonce[:nc.size]
}
