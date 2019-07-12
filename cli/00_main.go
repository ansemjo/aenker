// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

// Package cli implements the commandline interface using https://github.com/spf13/cobra.
package cli

import (
	"os"
	"path"

	"github.com/spf13/cobra"
)

// RootCommand is the root cobra command to be executed with Execute().
var RootCommand = &cobra.Command{
	Use: "aenker",
	Long: `Aenker is a tool to encrypt files with an authenticated integrated encryption
scheme by chunking the input into equal parts and sealing them with a key
derived from an anonymous Diffie-Hellman key exchange on an elliptic curve.

Many parties can encrypt files for a single recipient by distributing that
recipient's public key, while only the recipient can decrypt any of those files
afterwards.`,
}

// Initialize cobra commander and disable sorting. When added via a
// file's init() function, they will appear in lexicographical order
// of their respective files.
func init() {
	this := RootCommand
	cobra.EnableCommandSorting = false
	this.Flags().SortFlags = false
}

// Execute is the main function. It starts the cobra commander for the RootCommand 'aenker',
// parses arguments and flags, and finally executes the desired command.
func Execute() {
	if err := RootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}

// the default secret key used by several commands
var defaultkey = func() string {
	if home, err := os.UserHomeDir(); err == nil {
		return path.Join(home, ".local", "share", "aenker", "aenkerkey")
	} else {
		return path.Join("./", "aenkerkey") // fallback to current dir
	}
}()
