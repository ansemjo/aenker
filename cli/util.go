// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package cli

import (
	"bufio"
	b64 "encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/ansemjo/aenker/keyderivation"
	"golang.org/x/crypto/ssh/terminal"
)

var base64 = b64.StdEncoding.EncodeToString

// Treat any non-nil error as a fatal failure,
// print error to stderr and exit with nonzero status.
func fatal(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "FATAL:", err)
		os.Exit(1)
	}
}

// read password and derive key
func getpasskey(key *[32]byte, reader io.Reader) (err error) {

	var passwd []byte

	// try interactive if terminal
	stdin := int(os.Stdin.Fd())
	if terminal.IsTerminal(stdin) {

		fmt.Fprint(os.Stderr, "Enter password: ")
		passwd, err = terminal.ReadPassword(stdin)
		fmt.Fprint(os.Stderr, "\n")

	} else if reader != nil {

		buf := bufio.NewReader(os.Stdin)
		passwd, _, err = buf.ReadLine()

	} else {

		err = errors.New("cannot read password: stdin is not a terminal")

	}

	if err != nil {
		return
	}

	// derive key
	k := keyderivation.Password(passwd, "aenker")
	copy(key[:], k)

	return

}
