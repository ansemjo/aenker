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
	Use:     "manual [man|markdown]",
	Aliases: []string{"man", "documentation"},
	Short:   "Output documentation",
	Long:    "Generate documentation in manpage or markdown format.",
	Example: `
Put manual pages in your local manpath:
  aenker gen manual -d ~/.local/share/man/`,
	ValidArgs: []string{"man", "markdown"},
	Args:      cobra.OnlyValidArgs,
	PreRunE:   checkOutDirFlag,
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		// force cobra to use a simple name, otherwise it could
		// contain path seperator characters (e.g. './aenker')
		rootCmd.Use = docsname

		if len(args) > 0 && args[0] == "markdown" {

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
