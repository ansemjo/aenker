// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

// +build !nokeygen

package cli

import (
	"fmt"
	"os"

	cf "github.com/ansemjo/aenker/cli/cobraflags"
	"github.com/ansemjo/aenker/keyderivation"
	"github.com/spf13/cobra"
)

// AddPubkeyCommand adds the pubkey converter subcommand to a cobra command.
func AddPubkeyCommand(parent *cobra.Command) *cobra.Command {

	var private *cf.Key32Flag
	var public *cf.FileFlag

	command := &cobra.Command{
		Use:     "pubkey",
		Aliases: []string{"pk"},
		Short:   "calculate public key",
		Long:    "Calculate the public key of a Curve25519 private key.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return cf.CheckAll(cmd, args, public.Open, private.Check)
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			// close open file upon exit
			defer func() { public.File.Close() }()

			// calculate and write public part
			pub := keyderivation.Public(private.Key)
			_, err = fmt.Fprintln(public.File, base64(pub[:]))

			return
		},
	}
	command.Flags().SortFlags = false

	// add the output file flags
	private = cf.AddKey32Flag(command, "key", "k", "private key (default: stdin)", os.Stdin)
	public = cf.AddFileFlag(command, "pubkey", "p", "write public key to file (default: stdout)",
		cf.Exclusive(0644), os.Stdout)

	parent.AddCommand(command)
	return command
}