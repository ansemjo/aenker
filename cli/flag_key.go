// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package cli

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/spf13/cobra"
)

type CmdCheckFunc func(cmd *cobra.Command) error

type Key32Flag struct {
	Key   *[]byte
	Check CmdCheckFunc
}

type Key32FlagOptions struct {
	Flag         string
	Short        string
	Usage        string
	DefaultStdin bool
}

// AddKey32Flag adds a flag to a command, which can either be a valid base64
// string or a filename for a 32 byte key. Optionally reads from stdin.
func AddKey32Flag(cmd *cobra.Command, opts Key32FlagOptions) *Key32Flag {

	// add flag to command
	str := cmd.Flags().StringP(opts.Flag, opts.Short, "", opts.Usage)
	key := &Key32Flag{}

	// return struct and build check function inline
	key.Check = func(cmd *cobra.Command) (err error) {

		// if flag was given
		if cmd.Flag(opts.Flag).Changed {

			// and it is a valid base64 encoded key
			if is32ByteBase64Encoded(*str) {
				key.Key, err = decodeKey([]byte(*str))

			} else {
				// assume any other string to be a filename
				file, err := os.Open(*str)
				if err != nil {
					return err
				}
				defer file.Close()
				key.Key, err = decodeKeyFile(file)
			}

		} else if opts.DefaultStdin {
			// if flag was not given but "read from stdin" is true
			key.Key, err = decodeKeyFile(os.Stdin)
		}

		// if neither, just return nil. the pointer to Key will remain nil!
		return
	}
	return key
}

// is32ByteBase64Encoded checks if the given string is a base64-encoded 32 byte value.
func is32ByteBase64Encoded(str string) bool {
	return regexp.MustCompile("^[A-Za-z0-9+/]{43}=$").MatchString(str)
}

// decodeKey decodes a base64 string and expects a 32 byte value inside
func decodeKey(str []byte) (key *[]byte, err error) {
	k := make([]byte, 32)
	n, err := base64.StdEncoding.Decode(k, str)
	if err != nil {
		return
	}
	if n != 32 {
		err = errors.New("key must be 32 bytes")
		return
	}
	key = &k
	return
}

// decodeKeyFile reads a file and decodes its contents with decodeKey
func decodeKeyFile(file *os.File) (key *[]byte, err error) {
	keyslice, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}
	return decodeKey(keyslice)
}
