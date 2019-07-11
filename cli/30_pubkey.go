// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package cli

import (
	"fmt"
	"os"

	cf "github.com/ansemjo/aenker/cli/cobraflags"
	"github.com/ansemjo/aenker/keyderivation"
	"github.com/spf13/cobra"
)

func init() {
	AddPubkeyCommand(RootCommand)
}

// AddPubkeyCommand adds the pubkey converter subcommand to a cobra command.
func AddPubkeyCommand(parent *cobra.Command) *cobra.Command {

	var private *cf.Key32Flag

	command := &cobra.Command{
		Use:     "pubkey",
		Aliases: []string{"pk", "show"},
		Short:   "print public key",
		Long: `Calculate the public key of a Curve25519 private key by performing a base point
multiplication. You could use any source of 32 random bytes as input.
When called as "show" a formatted seal command will be shown.`,
		Example: `  # show default key
  aenker show

  # new keypair from system randomness
  head -c32 /dev/urandom | base64 > mykey
  aenker pk -k mykey > mykey.pub`,

		Args: cf.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return private.Check(cmd, args)
		},

		RunE: func(cmd *cobra.Command, args []string) (err error) {

			// calculate public key
			pub := base64(keyderivation.Public(private.Key)[:])

			// write formatted seal command if called as "show"
			if cmd.CalledAs() == "show" {
				_, err = fmt.Printf(
					"Encrypt files to %q with:\n\n"+
						"  aenker seal -p %s ...\n\n", private.File, pub)
			} else {
				_, err = fmt.Println(pub)
			}

			return
		},
	}
	command.Flags().SortFlags = false

	// add the input keyfile flag
	private = cf.AddKey32Flag(command, "key", "k", defaultkey, "private key (default: stdin)", os.Stdin)

	parent.AddCommand(command)
	return command
}
