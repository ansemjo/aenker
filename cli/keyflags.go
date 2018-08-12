// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package cli

import (
	"bufio"
	"encoding/base64"
	"errors"
	"os"

	"github.com/mitchellh/go-homedir"

	"github.com/spf13/cobra"
)

var keyfile string
var key []byte

// add necessary key flags to a command
func addKeyFlags(cmd *cobra.Command) {

	cmd.Flags().StringVarP(&keyfile, "keyfile", "f", "~"+string(os.PathSeparator)+".aenker", "file with the key on the first line")
	cmd.Flags().BytesBase64VarP(&key, "key", "k", nil, "encoded key as string")

}

// check and load keys .. run this in PreRunE
func checkKeyFlags(cmd *cobra.Command, args []string) (err error) {

	fileGiven := cmd.Flag("keyfile").Changed
	keyGiven := cmd.Flag("key").Changed

	if fileGiven && keyGiven { // both were given
		return errors.New("only use either one of keyfile or key")
	}

	if keyGiven && len(key) != 32 { // key was given and it's not 32 bytes
		return errors.New("key must be 32 bytes")
	}

	if !keyGiven { // key was not given, use (default) file

		// expand ~ only if this was the default, otherwise it might
		// be a literal tilde ..
		if !fileGiven {
			keyfile, err = homedir.Expand(keyfile)
			if err != nil {
				return
			}
		}

		f, err := os.Open(keyfile) // open keyfile for reading
		if err != nil {
			return err
		}
		defer f.Close()

		line, _, err := bufio.NewReader(f).ReadLine() // read the first line
		if err != nil {
			return err
		}

		n, err := base64.StdEncoding.Decode(line, line) // decode base64 in line
		if err != nil {
			return err
		}

		if n != 32 { // decoded slice needs to be 32 bytes
			return errors.New("key must be 32 bytes")
		}

		key = line[:n] // put it in the key slice
	}

	return nil
}
