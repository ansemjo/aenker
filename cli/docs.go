package cli

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// the command name to used in docs
var docsname string

func init() {
	docsCmd.Flags().SortFlags = false
	addOutDirFlag(docsCmd)
	docsCmd.MarkFlagRequired("dir")
	docsCmd.Flags().StringVarP(&docsname, "name", "n", "aenker", "name of this program to be used in manuals")
}

var docsCmd = &cobra.Command{
	Use:     "docs {man|markdown}",
	Aliases: []string{"documentation", "man", "manual"},
	Short:   "output manual / documentation",
	Long:    "Generate documentation in manpage or markdown format.",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		err = cobra.ExactArgs(1)(cmd, args)
		if err != nil {
			return
		}
		err = cobra.OnlyValidArgs(cmd, args)
		return
	},
	ValidArgs: []string{"man", "markdown"},
	PreRunE:   checkOutDirFlag,
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		// force cobra to use a simple name, otherwise it could
		// contain path seperator characters (e.g. './aenker')
		rootCmd.Use = docsname

		if args[0] == "markdown" {

			err = doc.GenMarkdownTree(rootCmd, outdir)

		} else {

			// manpages need to be in subdirs manN, where N is the section
			subdir := outdir + string(os.PathSeparator) + "man1"
			err = os.MkdirAll(subdir, 0755)
			if err != nil {
				return
			}

			header := &doc.GenManHeader{Title: "AENKER", Section: "1"}
			err = doc.GenManTree(rootCmd, header, subdir)

		}

		return
	},
}
