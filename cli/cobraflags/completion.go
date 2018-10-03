// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package cobraflags

import (
	"os"

	"github.com/spf13/cobra"
)

func AddCompletionCommand(parent, root *cobra.Command) *cobra.Command {

	var file *FileFlag

	command := &cobra.Command{
		Use:       "completion [bash|zsh]",
		Short:     "generate autocompletion",
		Long:      "Generate autocompletion scripts to be sourced by your shell.",
		ValidArgs: []string{"bash", "zsh"},
		Args: func(cmd *cobra.Command, args []string) error {
			return CheckAll(cmd, args, cobra.MaximumNArgs(1), cobra.OnlyValidArgs)
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return CheckAll(cmd, args, file.Open)
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			// generate bash completions, unless zsh is given
			if len(args) > 0 && args[0] == "zsh" {
				err = root.GenZshCompletion(file.File)
			} else {
				err = root.GenBashCompletion(file.File)
			}
			file.File.Close()
			return

		},
	}

	file = AddFileFlag(command,
		"out", "o", "output completion script to file",
		func(name string) (*os.File, error) {
			return os.OpenFile(name, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		},
		os.Stdout)

	parent.AddCommand(command)
	return command
}
