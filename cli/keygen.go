package cli

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(keygen)
}

var keygen = &cobra.Command{
	Use:   "kg",
	Short: "Generate a random 32-byte key",
	Long:  "Generate a random 32-byte key and output base64-encoded form to stdout",
	Run: func(cmd *cobra.Command, args []string) {

		r := make([]byte, 32)
		rand.Read(r)

		fmt.Println(base64.StdEncoding.EncodeToString(r))

	},
}
