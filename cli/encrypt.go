package cli

import (
	"fmt"
	"os"

	"github.com/ansemjo/aenker/aenker"
	"github.com/spf13/cobra"
)

func init() {
	root.AddCommand(encryptCmd)
	encryptCmd.Flags().StringVarP(&outputFile, "output", "o", "", "optional file to write to")
}

var outputFile string

var encryptCmd = &cobra.Command{
	Use:   "enc",
	Short: "encrypt data with aenker",
	Long:  "encrypt stdin and place the ciphertext in stdout",
	Run:   encrypt,
}

func encrypt(cmd *cobra.Command, args []string) {

	zk := make([]byte, aenker.KeyLength)
	ae, _ := aenker.NewAenker(zk)

	lw, err := ae.Encrypt(os.Stdout, os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "wrote %d bytes\n", lw)

}
