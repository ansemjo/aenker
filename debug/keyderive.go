package main

import (
	"encoding/base64"
	"fmt"

	"github.com/ansemjo/aenker/ae/keyderivation"
)

func keyderive() {

	secret := []byte("pass")
	salt := []byte("salt")

	b64 := base64.StdEncoding.EncodeToString

	fmt.Println("symmetric :", b64(keyderivation.Symmetric(secret, salt)))
	fmt.Println("password  :", b64(keyderivation.Password(secret, salt)))
	// Elliptic accepts slices < 32 bytes but will pad them with zeros internally
	// using a nil salt produces the same key regardless of secret though, careful!
	fmt.Println("elliptic  :", b64(keyderivation.Elliptic(secret, salt)))

}
