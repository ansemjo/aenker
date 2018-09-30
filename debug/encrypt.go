package main

import (
	"os"
	"strings"

	"github.com/ansemjo/aenker/ae"
)

func encrypt() {

	key := make([]byte, 32)
	data := strings.NewReader("Hello, World!")

	ae := ae.NewAenker(&key)
	ae.Encrypt(os.Stdout, data, 8)

}
