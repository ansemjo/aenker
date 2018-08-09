package cli

import (
	"fmt"
	"os"

	"github.com/ansemjo/aenker/aenker"
	"github.com/spf13/cobra"
)

func init() {
	root.AddCommand(decryptCmd)
}

var decryptCmd = &cobra.Command{
	Use:   "dec",
	Short: "decrypt data with aenker",
	Long:  "decrypt stdin and place the plaintext in stdout",
	Run:   decrypt,
}

func decrypt(cmd *cobra.Command, args []string) {

	zk := make([]byte, aenker.KeyLength)
	ae, _ := aenker.NewAenker(zk)

	lw, err := ae.Decrypt(os.Stdout, os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "wrote %d bytes\n", lw)

}
