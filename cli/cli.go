package cli

import (
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
	// Run: func(cmd *cobra.Command, args []string) {
	// 	cmd.GenBashCompletion(os.Stdout)
	// },
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
