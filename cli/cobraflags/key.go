// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package cobraflags

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"
)

type Key32Flag struct {
	Key   *[32]byte
	Check func(cmd *cobra.Command, args []string) error
}

// AddKey32Flag adds a flag to a command, which can either be a valid base64
// string or a filename for a 32 byte key. Optionally reads from stdin.
func AddKey32Flag(cmd *cobra.Command, flag, short, usage string, fallback *os.File) (kf *Key32Flag) {

	// add flag to command
	str := cmd.Flags().StringP(flag, short, "", usage)

	// return struct with check function for PreRunE
	return &Key32Flag{
		Check: func(cmd *cobra.Command, args []string) (err error) {
			if cmd.Flag(flag).Changed {

				// given string is a valid key
				if is32ByteBase64Encoded(*str) {
					kf.Key, err = decodeKey(*str)

				} else {
					// assume any other string to be a filename
					var file *os.File
					file, err = os.Open(*str)
					if err != nil {
						return err
					}
					defer file.Close()
					kf.Key, err = decodeKeyFile(file)
				}

			} else if fallback != nil {
				// if flag was not given but a fallback was defined
				kf.Key, err = decodeKeyFile(fallback)
			}

			// if neither Key will remain nil!
			return
		},
	}
}

// is32ByteBase64Encoded checks if the given string is a base64-encoded 32 byte value.
func is32ByteBase64Encoded(str string) bool {
	return regexp.MustCompile("^[A-Za-z0-9+/]{43}=$").MatchString(str)
}

// decodeKey decodes a base64 string and expects a 32 byte value inside
func decodeKey(str string) (key *[32]byte, err error) {

	k, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return
	}

	if len(k) != 32 {
		err = errors.New("key must be 32 bytes")
		return
	}

	key = new([32]byte)
	copy(key[:], k)
	return
}

// decodeKeyFile reads a file and decodes its contents with decodeKey
func decodeKeyFile(file *os.File) (key *[32]byte, err error) {

	// use a line scanner
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// test each line for key regexp
		if line := scanner.Text(); is32ByteBase64Encoded(line) {
			return decodeKey(line)
		}
	}

	// return any errors encountered
	if e := scanner.Err(); e != nil {
		return nil, e
	}

	// probably hit EOF
	return nil, fmt.Errorf("no base64 encoded key found in %s", file.Name())

}
