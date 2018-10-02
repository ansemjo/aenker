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

func EncryptCommand(parent *cobra.Command) {

	var key *cf.Key32Flag
	command := &cobra.Command{

		Use:     "encrypt",
		Aliases: []string{"enc", "e"},
		Short:   "encrypt a file",
		Long:    "Encrypt Stdin and write the ciphertext to Stdout.",

		PreRunE: func(cmd *cobra.Command, args []string) (err error) {
			return key.Check(cmd)
		},

		RunE: func(cmd *cobra.Command, args []string) (err error) {

			ae, err := ae.NewWriter(os.Stdout, key.Key)
			fatal(err)
			defer ae.Close()
			io.Copy(ae, os.Stdin)
			return

		},
	}
	command.Flags().SortFlags = false

	// add required peer key flag
	key = cf.AddKey32Flag(command, "peer", "p", "receiver's public key", nil)
	command.MarkFlagRequired("key")

	parent.AddCommand(command)
}
