// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package ae

import (
	"crypto/cipher"
)

type kdfMode byte

const (
	kdfModeSymmetric = 0
	kdfModePassword  = 1
	kdfModeElliptic  = 2
	kdfModeMax       = kdfModeElliptic
)

type Aenker2 struct {
	aead      cipher.AEAD
	kdf       kdfMode
	chunksize uint16
	nonce     []byte
}
