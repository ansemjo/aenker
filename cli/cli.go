// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   os.Args[0],
	Short: "authenticated encryption on the commandline",
	Long: `aenker is a tool to encrypt files with an authenticated
cipher (ChaCha20Poly1305) in a 'streamable' way by chunking
the input into equally-sized parts.`,
	Version: "0.3.5 (not built with build.go)",
	Example: `
Generate a new random key:
  aenker kg -o ~/.aenker

Encrypt a file:
  aenker enc -i /path/to/secret/documents.tar.gz -o encrypted.tar.gz.ae

Encrypt using pipes and redirection:
  echo 'Hello, World!' | aenker e -f ./otherkey > hello.ae

Decrypt and unpack an encrypted tar archive:
  aenker dec -i encrypted.tar.gz.ae | tar -xzv`,
}

var gendocCmd = &cobra.Command{
	Use:     "gen",
	Aliases: []string{"generate"},
	Short:   "Generate documentation or autocompletion",
}

// Initialize cobra commander, disable sorting and
// add commands. We do that here instead of in individual
// file init()s because we want to define the sorting manually.
func init() {
	cobra.EnableCommandSorting = false
	rootCmd.Flags().SortFlags = false
	rootCmd.AddCommand(encryptCmd)
	rootCmd.AddCommand(decryptCmd)
	rootCmd.AddCommand(keygenCmd)
	rootCmd.AddCommand(gendocCmd)
	gendocCmd.AddCommand(docsCmd)
	gendocCmd.AddCommand(completionCmd)
}

// SetVersion sets the version string if a more specific
// or updated string is known
func SetVersion(version string) {
	rootCmd.Version = version
}

// Execute is the main function. It starts the cobra commander,
// parses arguments and flags, and finally executes the desired command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// Treat any non-nil error as a fatal failure,
// print error to stderr and exit with nonzero status.
func fatal(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
