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

		Use:     "seal",
		Aliases: []string{"encrypt", "e"},
		Short:   "encrypt a file",
		Long:    "Encrypt a file for a recipient's public key and output authenticated ciphertext.",
		Example: "  tar -cz * | aenker seal -p $PUBLICKEY > archive.tar.gz.ae",

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return cf.CheckAll(cmd, args, key.Check, input.Open, output.Open)
		},

		Run: func(cmd *cobra.Command, args []string) {

			ae, err := ae.NewWriter(output.File, key.Key)
			fatal(err)
			defer ae.Close()

			_, err = io.Copy(ae, input.File)
			fatal(err)

			return

		},
	}
	command.Flags().SortFlags = false

	// add required peer key flag
	key = cf.AddKey32Flag(command, "peer", "p", "", "receiver's public key", nil)
	command.MarkFlagRequired("peer")

	// add input/output flags
	input = cf.AddFileFlag(command, "input", "i", "input file, plaintext (default: stdin)",
		cf.Readonly(), os.Stdin)
	output = cf.AddFileFlag(command, "output", "o", "output file, ciphertext (default: stdout)",
		cf.Truncate(0644), os.Stdout)

	parent.AddCommand(command)
	return command
}
