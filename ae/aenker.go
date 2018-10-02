// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package ae

type header struct {
	magic     [8]byte
	salt      [8]byte
	ephemeral [32]byte
}
