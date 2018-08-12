// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package aenker

import "encoding/binary"

// Nonce is a counter which is monotonically incremented
// and outputs a 12-byte slice on .Next() for AEAD usage
type nonceCounter struct {
	ctr uint64
}

// NewNonce returns a new Nonce starting at 0
func newNonceCounter() *nonceCounter {
	return &nonceCounter{}
}

// Next outputs the current counter value as a 12-byte
// slice and increments the internal counter
func (nc *nonceCounter) Next() (nonce []byte) {
	nonce = make([]byte, 12)
	binary.LittleEndian.PutUint64(nonce, nc.ctr)
	nc.ctr++
	return
}
