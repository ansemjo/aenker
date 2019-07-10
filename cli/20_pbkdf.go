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

	cf "github.com/ansemjo/aenker/cli/cobraflags"
	"github.com/ansemjo/aenker/keyderivation"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

func init() {

	// AddPbkdfCommand adds the password-based key generator subcommand to a cobra command.
	//
	// It must be explicitly enabled with `-tags pbkdf`. You can use https://github.com/ansemjo/stdkdf instead.
	AddPbkdfCommand = func(parent *cobra.Command) *cobra.Command {

		var private *cf.FileFlag
		var public *cf.FileFlag
		var salt string

		command := &cobra.Command{
			Use:   "pbkdf",
			Short: "generate a password-derived keypair",
			Long: `Generate and save a new Curve25519 keypair by deriving it from a password with
Argon2i. The cost settings are predefined as time=32, memory=256MB, threads=4.`,
			Example: "  aenker kg pbkdf -s mysaltstring -o mykey -p mykey.pub",
			PreRunE: func(cmd *cobra.Command, args []string) (err error) {

				// output file
				err = private.Open(cmd, args)
				if err != nil {
					return
				}

				// public key file
				err = public.Open(cmd, args)
				if err != nil {
					os.Remove(private.File.Name())
					return
				}

				return
			},
			RunE: func(cmd *cobra.Command, args []string) (err error) {

				// close all open file upon exit
				defer func() {
					for _, f := range []*cf.FileFlag{private, public} {
						if f.File != nil {
							f.File.Close()
						}
					}
				}()

				// derive key from password
				key := new([32]byte)
				err = getpasskey(key, salt, os.Stdin)
				fatal(err)

				// write encoded key to file
				if private.File == os.Stdout {
					fmt.Fprintf(os.Stderr, "private key:\n  ")
				}
				_, err = fmt.Fprintln(private.File, base64(key[:]))
				fatal(err)

				// if public was given, write public part
				if public.File != nil {

					pub := keyderivation.Public(key)
					if public.File == os.Stdout {
						fmt.Fprintf(os.Stderr, "public key:\n  ")
					}
					_, err = fmt.Fprintln(public.File, base64(pub[:]))

				}

				return
			},
		}
		command.Flags().SortFlags = false

		// add the output file flags
		private = cf.AddFileFlag(command, "out", "o", "write output to file (default: stdout)",
			cf.Truncate(0600), os.Stdout)

		public = cf.AddFileFlag(command, "pubkey", "p", "write public key to file (default: stdout)",
			cf.Truncate(0644), os.Stdout)

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
