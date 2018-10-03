// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package cobraflags

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// -------------------------- generate --------------------------

func AddGeneratorCommand(parent *cobra.Command) *cobra.Command {

	command := &cobra.Command{
		Use:   "docs",
		Short: "generate documentation",
		Long:  "Generate documentation or autocompletion scripts.",
	}

	// add content to be generated
	AddCompletionCommand(command, parent)
	AddDocumentationCommand(command, parent)

	parent.AddCommand(command)
	return command
}

// -------------------------- documentation --------------------------

func AddDocumentationCommand(parent, root *cobra.Command) *cobra.Command {

	var dir *DirFlag

	command := &cobra.Command{
		Use:       "manual [man|markdown]",
		Aliases:   []string{"man", "docs"},
		Short:     "generate manuals",
		Long:      "Generate documentation in manpage or markdown format.",
		ValidArgs: []string{"man", "markdown"},
		Args: func(cmd *cobra.Command, args []string) error {
			return CheckAll(cmd, args, cobra.MaximumNArgs(1), cobra.OnlyValidArgs)
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return CheckAll(cmd, args, dir.Check)
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			if len(args) > 0 && args[0] == "markdown" {

				err = doc.GenMarkdownTree(root, dir.Dir)

			} else {

				// manpages need to be in subdirs manN, where N is the section
				subdir := dir.Dir + string(os.PathSeparator) + "man1"
				err = os.MkdirAll(subdir, 0755)
				if err != nil {
					return
				}

				header := &doc.GenManHeader{Title: strings.ToUpper(root.Name()), Section: "1"}
				err = doc.GenManTree(root, header, subdir)

			}

			return
		},
	}

	// add output directory flag
	dir = AddDirFlag(command, "directory", "d", "", "output directory")
	command.MarkFlagRequired("directory")

	parent.AddCommand(command)
	return command
}
