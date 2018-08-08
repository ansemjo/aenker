package main

import "encoding/binary"

type NonceCounter struct {
	ctr uint64
}

func NewNonceCounter() *NonceCounter {
	return &NonceCounter{}
}

func (nc *NonceCounter) Next() (nonce []byte) {
	nonce = make([]byte, 12)
	binary.LittleEndian.PutUint64(nonce, nc.ctr)
	nc.ctr++
	return
}
