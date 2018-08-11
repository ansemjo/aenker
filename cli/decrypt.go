package cli

import (
	"os"

	"github.com/ansemjo/aenker/Aenker"
	"github.com/spf13/cobra"
)

func init() {
	decryptCmd.Flags().SortFlags = false
	addKeyFlags(decryptCmd)
}

var decryptCmd = &cobra.Command{
	Use:   "d",
	Short: "decrypt a file",
	Long:  "decrypt stdin and place the plaintext in stdout",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {

		return checkKeyFlags(cmd, args)

	},
	Run: func(cmd *cobra.Command, args []string) {

		reader := os.Stdin
		writer := os.Stdout

		ae := aenker.NewAenker(&key)
		_, err := ae.Decrypt(writer, reader)
		fatal(err)

	},
}
