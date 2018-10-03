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
	command := &cobra.Command{

		Use:     "decrypt",
		Aliases: []string{"open", "d"},
		Short:   "decrypt a file",
		Long:    "Decrypt from Stdin and write the plaintext to Stdout.",

		PreRunE: func(cmd *cobra.Command, args []string) (err error) {
			return key.Check(cmd)
		},

		Run: func(cmd *cobra.Command, args []string) {
			var err error
			defer func() { fatal(err) }()

			ae, err := ae.NewReader(os.Stdin, key.Key)
			fatal(err)
			_, err = io.Copy(os.Stdout, ae)
			return

		},
	}
	command.Flags().SortFlags = false

	// add required private key flag
	key = cf.AddKey32Flag(command, "key", "k", "your private key", nil)
	command.MarkFlagRequired("key")

	parent.AddCommand(command)
	return command
}
