package aenker

import "encoding/binary"

// Nonce is a counter which is monotonically incremented
// and outputs a 12-byte slice on .Next() for AEAD usage
type Nonce struct {
	ctr uint64
}

// NewNonce returns a new Nonce starting at 0
func NewNonce() *Nonce {
	return &Nonce{}
}

// Next outputs the current counter value as a 12-byte
// slice and increments the internal counter
func (nc *Nonce) Next() (nonce []byte) {
	nonce = make([]byte, 12)
	binary.LittleEndian.PutUint64(nonce, nc.ctr)
	nc.ctr++
	return
}
