package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "aenker",
	Short: "aenker is an authenticated encryptor",
	Long: `aenker is a tool that operates an AEAD
(ChaCha20Poly1305) in a streamable way
by chunking the input in equally-sized
parts.`,
	Version: "none",
}

func init() {
	cobra.EnableCommandSorting = false
	rootCmd.Flags().SortFlags = false
	// add commands here instead of in individual file inits
	// because we want to define the sorting ourselves
	rootCmd.AddCommand(encryptCmd)
	rootCmd.AddCommand(decryptCmd)
	rootCmd.AddCommand(keygenCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// any non-nil error is a fatal failure.
// print error to stderr and exit
func fatal(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
