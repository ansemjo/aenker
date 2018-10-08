// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

// Package cli implements the commandline interface using https://github.com/spf13/cobra.
package cli

import (
	"os"

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
	Version: SpecificVersion(),
}

// Initialize cobra commander, disable sorting and
// add commands. We do that here instead of in individual
// file init()s because we want to define the sorting manually.
func init() {
	this := RootCommand
	cobra.EnableCommandSorting = false
	this.Flags().SortFlags = false

	AddEncryptCommand(this)
	AddDecryptCommand(this)
}

// Execute is the main function. It starts the cobra commander for the RootCommand 'aenker',
// parses arguments and flags, and finally executes the desired command.
func Execute() {
	if err := RootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
