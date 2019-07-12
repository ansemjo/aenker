// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

// +build !nokeygen,pbkdf

package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	cf "github.com/ansemjo/aenker/cli/cobraflags"
	"github.com/ansemjo/aenker/keyderivation"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

// used to restore state upon abnormal exit
var initialState *terminal.State

func init() {

	// from: https://groups.google.com/d/msg/golang-nuts/kTVAbtee9UA/Y1F5MbASCQAJ
	// remember initial terminal state
	var err error
	if initialState, err = terminal.GetState(syscall.Stdin); err != nil {
		return
	}
	// and restore it on exit
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		<-c
		_ = terminal.Restore(syscall.Stdin, initialState)
		fmt.Print("\n")
		os.Exit(0)
	}()

	// AddPbkdfCommand adds the password-based key generator subcommand to a cobra command.
	// It must be explicitly enabled with `-tags pbkdf` or you can use
	// https://github.com/ansemjo/stdkdf instead.
	AddPbkdfCommand = func(parent *cobra.Command) *cobra.Command {

		var keyfile, salt string

		command := &cobra.Command{
			Use:   "pbkdf",
			Short: "generate a password-derived keypair",
			Long: `Generate and save a new Curve25519 keypair by deriving it from a password with
Argon2i. The cost settings are predefined as time=32, memory=256MB, threads=4.`,
			Example: "  aenker kg pbkdf -s mysaltstring -o mykey -p mykey.pub",
			Args:    cf.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) (err error) {

				// format returned errors
				defer func() {
					if err != nil {
						err = fmt.Errorf("aenker keygen: %s", err)
						fatal(err)
					}
				}()

				// derive key from password
				seckey := new([32]byte)
				if err = getpasskey(seckey, salt, os.Stdin); err != nil {
					return
				}

				// write to file and return pubkey
				pubkey, err := writeKey(seckey, keyfile,
					fmt.Sprintf("version: %s, salt: %s", RootCommand.Version, salt))
				if err != nil {
					return
				}

				// print info to stdout
				fmt.Printf(`Derived key saved in %q.
Your public key is: %s
Use the following command to encrypt files for this key:

  aenker seal -p %s ...

`, keyfile, pubkey, pubkey)

				return

			},
		}
		command.Flags().SortFlags = false

		// add the output file flag
		command.Flags().StringVarP(&keyfile, "file", "f", defaultkey, "save key to this file")

		// add the salt flag
		command.Flags().StringVarP(&salt, "salt", "s", "aenker", "salt for argon2i key derivation")

		parent.AddCommand(command)
		return command
	}

}

// read password and derive key
func getpasskey(key *[32]byte, salt string, reader io.Reader) (err error) {

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
	k := keyderivation.Password(passwd, salt)
	copy(key[:], k)

	return

}
