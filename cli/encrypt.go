package cli

import (
	"fmt"
	"os"

	"github.com/ansemjo/aenker/aenker"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(encryptCmd)
	encryptCmd.Flags().SortFlags = false
	addKeyFlags(encryptCmd)
	addChunkSizeFlag(encryptCmd)

}

var encryptCmd = &cobra.Command{
	Use:   "e",
	Short: "encrypt a file",
	Long:  "encrypt stdin and place the ciphertext in stdout",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {

		err = parseChunkSize(cmd, args)
		if err != nil {
			return
		}

		return checkKeyFlags(cmd, args)

	},
	Run: func(cmd *cobra.Command, args []string) {

		reader := os.Stdin
		writer := os.Stdout

		mek, blob, err := aenker.NewMEK(key)
		fatal(err)

		ae, err := aenker.NewAenker(mek, chunksize)
		fatal(err)

		_, err = os.Stdout.Write(blob)
		fatal(err)

		lw, err := ae.Encrypt(writer, reader)
		fatal(err)

		fmt.Fprintf(os.Stderr, "wrote %d bytes\n", lw)

	},
}
