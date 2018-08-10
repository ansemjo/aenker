package cli

import (
	"fmt"
	"os"

	"github.com/ansemjo/aenker/aenker"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(decryptCmd)
	decryptCmd.Flags().SortFlags = false
	addKeyFlags(decryptCmd)
	addChunkSizeFlag(decryptCmd)
}

var decryptCmd = &cobra.Command{
	Use:   "d",
	Short: "decrypt a file",
	Long:  "decrypt stdin and place the plaintext in stdout",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		err = parseChunkSize(cmd)
		if err != nil {
			return
		}
		return checkKeyFlags(cmd)
	},
	Run: decrypt,
}

func decrypt(cmd *cobra.Command, args []string) {

	ae, _ := aenker.NewAenker(key, chunksize)

	lw, err := ae.Decrypt(os.Stdout, os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "wrote %d bytes\n", lw)

}
