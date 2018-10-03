// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package cli

import (
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

	command := &cobra.Command{

		Use:     "decrypt",
		Aliases: []string{"open", "d"},
		Short:   "decrypt a file",
		Long:    "Decrypt from Stdin and write the plaintext to Stdout.",

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return cf.CheckAll(cmd, args, key.Check, input.Open, output.Open)
		},

		Run: func(cmd *cobra.Command, args []string) {
			var err error
			defer func() { fatal(err) }()

			ae, err := ae.NewReader(input.File, key.Key)
			fatal(err)
			_, err = io.Copy(output.File, ae)
			return

		},
	}
	command.Flags().SortFlags = false

	// add required private key flag
	key = cf.AddKey32Flag(command, "key", "k", "your private key", nil)
	command.MarkFlagRequired("key")

	// add input/output flags
	input = cf.AddFileFlag(command, "input", "i", "input file, ciphertext (default: stdin)",
		cf.Readonly(), os.Stdin)
	output = cf.AddFileFlag(command, "output", "o", "output file, plaintext (default: stdout)",
		cf.Truncate(0644), os.Stdout)

	parent.AddCommand(command)
	return command
}
