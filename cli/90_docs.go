// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

// +build !nodocs

package cli

import (
	cf "github.com/ansemjo/aenker/cli/cobraflags"
	"github.com/spf13/cobra"
)

func init() {
	AddDocsGenCommand(RootCommand)
}

// AddDocsGenCommand adds the manuals and autocompletion generator subcommands to
// a cobra command.
//
// It can be disabled by building with the tag 'nodocs' to save some space.
func AddDocsGenCommand(root *cobra.Command) *cobra.Command {
	return cf.AddGeneratorCommand(root)
}
