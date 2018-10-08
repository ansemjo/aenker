// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package cli

import (
	"errors"
	"io"
	"os"

	"github.com/ansemjo/aenker/ae"
	cf "github.com/ansemjo/aenker/cli/cobraflags"
	"github.com/spf13/cobra"
)

// AddDecryptCommand adds the decryption subcommand to a cobra command.
func AddDecryptCommand(parent *cobra.Command) *cobra.Command {

	var key *cf.Key32Flag

	var input *cf.FileFlag
	var output *cf.FileFlag

	var password bool
	var salt string

	command := &cobra.Command{

		Use:     "open",
		Aliases: []string{"decrypt", "d"},
		Short:   "decrypt a file",
		Long:    "Decrypt from Stdin and write the plaintext to Stdout.",

		PreRunE: func(cmd *cobra.Command, args []string) (err error) {

			err = cf.CheckAll(cmd, args, key.Check, input.Open, output.Open)
			if err != nil {
				return
			}

			ch := func(n string) bool {
				return cmd.Flag(n).Changed
			}
			if !(ch("key") || ch("password")) || (ch("key") && ch("password")) {
				err = errors.New("either key or password required")
			}

			return
		},

		Run: func(cmd *cobra.Command, args []string) {

			if password {
				key.Key = new([32]byte)
				err := getpasskey(key.Key, salt, nil)
				fatal(err)
			}

			ae, err := ae.NewReader(input.File, key.Key)
			fatal(err)

			_, err = io.Copy(output.File, ae)
			fatal(err)

			return

		},
	}
	command.Flags().SortFlags = false

	// add required private key flag
	key = cf.AddKey32Flag(command, "key", "k", "your private key", nil)

	// add password flag
	command.Flags().BoolVarP(&password, "password", "p", false, "derive private key from password")
	command.Flags().StringVarP(&salt, "salt", "s", "aenker", "salt for password-based key derivation")

	// add input/output flags
	input = cf.AddFileFlag(command, "input", "i", "input file, ciphertext (default: stdin)",
		cf.Readonly(), os.Stdin)
	output = cf.AddFileFlag(command, "output", "o", "output file, plaintext (default: stdout)",
		cf.Truncate(0644), os.Stdout)

	parent.AddCommand(command)
	return command
}
