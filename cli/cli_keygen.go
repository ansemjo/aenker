// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package cli

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	cf "github.com/ansemjo/aenker/cli/cobraflags"
	"github.com/spf13/cobra"
)

func init() {
	this := keygenCmd
	this.Flags().SortFlags = false

	keygenSecKey = cf.AddFileFlag(this, cf.FileFlagOptions{
		"out", "o", "write output to file",
		func(name string) (*os.File, error) {
			return os.OpenFile(name, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		}, os.Stdout})

	keygenPubKey = cf.AddFileFlag(this, cf.FileFlagOptions{
		"pubkey", "p", "write public key to file",
		func(name string) (*os.File, error) {
			return os.OpenFile(name, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		}, nil})
}

var keygenSecKey *cf.FileFlag
var keygenPubKey *cf.FileFlag

var keygenCmd = &cobra.Command{
	Use:     "kg",
	Aliases: []string{"keygen"},
	Short:   "Generate a new key",
	Long:    "Generate a Curve25519 private key.",
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		// open output file
		err = keygenSecKey.Open(cmd)
		if err != nil {
			return
		}
		defer keygenSecKey.File.Close()

		// generate a new key and write encoded form to file
		key := NewBase64Key()
		_, err = fmt.Fprintln(keygenSecKey.File, key)
		fatal(err)

		return
	},
}

func NewBase64Key() string {
	key := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(key)
}
