package cli

import (
	"encoding/base64"
	"fmt"

	"github.com/ansemjo/aenker/aenker"
	"github.com/spf13/cobra"
)

var keygenCmd = &cobra.Command{
	Use:   "kg",
	Short: "generate a new key",
	Long:  "Generate a random 32-byte key and output base64-encoded form to stdout",
	Run: func(cmd *cobra.Command, args []string) {

		key := aenker.NewKey()
		fmt.Println(base64.StdEncoding.EncodeToString(key))

	},
}
