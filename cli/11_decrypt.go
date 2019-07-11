// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/ansemjo/aenker/ae"
	cf "github.com/ansemjo/aenker/cli/cobraflags"
	"github.com/spf13/cobra"
)

func init() {
	AddDecryptCommand(RootCommand)
}

// AddDecryptCommand adds the decryption subcommand to a cobra command.
func AddDecryptCommand(parent *cobra.Command) *cobra.Command {

	var key *cf.Key32Flag
	var input *cf.FileFlag
	var output *cf.FileFlag

	command := &cobra.Command{

		Use:     "open",
		Aliases: []string{"decrypt", "d"},
		Short:   "decrypt and authenticate a file",
		Long:    "Decrypt a file and output authenticated plaintext.",
		Example: "  aenker open -i archive.tar.gz.ae | tar -xz",

		PreRunE: func(cmd *cobra.Command, args []string) (err error) {

			// check in/out flags
			if err = cf.CheckAll(cmd, args, input.Open, output.Open); err != nil {
				return
			}

			// check key flag
			if err = key.Check(cmd, args); err != nil {
				err = fmt.Errorf("key is required: %s", err)
			}

			return
		},

		Run: func(cmd *cobra.Command, args []string) {

			ae, err := ae.NewReader(input.File, key.Key)
			fatal(err)

			_, err = io.Copy(output.File, ae)
			fatal(err)

			return

		},
	}
	command.Flags().SortFlags = false

	// add required private key flag
	key = cf.AddKey32Flag(command, "key", "k", defaultkey, "your private key", nil)

	// add input/output flags
	input = cf.AddFileFlag(command, "input", "i", "input file, ciphertext (default: stdin)",
		cf.Readonly(), os.Stdin)
	output = cf.AddFileFlag(command, "output", "o", "output file, plaintext (default: stdout)",
		cf.Truncate(0644), os.Stdout)

	parent.AddCommand(command)
	return command
}
