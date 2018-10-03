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

// AddEncryptCommand adds the encryption subcommand to a cobra command.
func AddEncryptCommand(parent *cobra.Command) *cobra.Command {

	var key *cf.Key32Flag

	var input *cf.FileFlag
	var output *cf.FileFlag

	command := &cobra.Command{

		Use:     "encrypt",
		Aliases: []string{"seal", "e"},
		Short:   "encrypt a file",
		Long:    "Encrypt Stdin and write the ciphertext to Stdout.",

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return cf.CheckAll(cmd, args, key.Check, input.Open, output.Open)
		},

		RunE: func(cmd *cobra.Command, args []string) (err error) {

			ae, err := ae.NewWriter(output.File, key.Key)
			fatal(err)
			defer ae.Close()
			io.Copy(ae, input.File)
			return

		},
	}
	command.Flags().SortFlags = false

	// add required peer key flag
	key = cf.AddKey32Flag(command, "peer", "p", "receiver's public key", nil)
	command.MarkFlagRequired("key")

	// add input/output flags
	input = cf.AddFileFlag(command, "input", "i", "input file, plaintext (default: stdin)",
		cf.Readonly(), os.Stdin)
	output = cf.AddFileFlag(command, "output", "o", "output file, ciphertext (default: stdout)",
		cf.Truncate(0644), os.Stdout)

	parent.AddCommand(command)
	return command
}
