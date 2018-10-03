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
	Long: `aenker is a tool to encrypt files with an authenticated
cipher (ChaCha20Poly1305) in a 'streamable' way by chunking
the input into equally-sized parts.`,
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
