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
	Run: encrypt,
}

func encrypt(cmd *cobra.Command, args []string) {

	ae, _ := aenker.NewAenker(key, chunksize)

	if cmd.Flag("chunksize").Changed {
		fmt.Fprintln(os.Stderr, "requested chunksize:", chunksize)
	}

	lw, err := ae.Encrypt(os.Stdout, os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "wrote %d bytes\n", lw)

}
