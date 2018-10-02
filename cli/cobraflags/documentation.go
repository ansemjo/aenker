// // Copyright (c) 2018 Anton Semjonov
// // Licensed under the MIT License

// TODO

package cobraflags

// import (
// 	"os"

// 	"github.com/spf13/cobra"
// 	"github.com/spf13/cobra/doc"
// )

// // -------------------------- generate --------------------------

// var generate = &cobra.Command{
// 	Use:     "gen",
// 	Aliases: []string{"generate"},
// 	Short:   "Generate documentation or autocompletion scripts.",
// }

// // -------------------------- documentation --------------------------

// var documentation = &cobra.Command{
// 	Use:       "manual [man|markdown]",
// 	Short:     "Generate documentation.",
// 	Long:      "Generate documentation in manpage or markdown format.",
// 	ValidArgs: []string{"man", "markdown"},
// 	Args:      cobra.OnlyValidArgs,
// 	PreRunE:   checkOutDirFlag,
// 	RunE: func(cmd *cobra.Command, args []string) (err error) {

// 		if len(args) > 0 && args[0] == "markdown" {

// 			err = doc.GenMarkdownTree(rootCmd, outdir)

// 		} else {

// 			// manpages need to be in subdirs manN, where N is the section
// 			subdir := outdir + string(os.PathSeparator) + "man1"
// 			err = os.MkdirAll(subdir, 0755)
// 			if err != nil {
// 				return
// 			}

// 			header := &doc.GenManHeader{Title: "AENKER", Section: "1"}
// 			err = doc.GenManTree(rootCmd, header, subdir)

// 		}

// 		return
// 	},
// }

// func init() {
// 	documentation.Flags().SortFlags = false
// 	addOutDirFlag(documentation)
// 	documentation.MarkFlagRequired("dir")
// 	completionCmd.Flags().SortFlags = false
// 	addOutFileFlag(completionCmd)
// }

// // -------------------------- completion --------------------------

// var completion *cobra.Command
// var completionFile *FileFlag

// func AddCompletionCommand(cmd *cobra.Command, rootcmd *cobra.Command) {
// 	completion = &cobra.Command{
// 		Use:       "completion [bash|zsh]",
// 		Short:     "Generate autocompletion scripts.",
// 		Long:      "Generate autocompletion scripts to be sourced by your shell.",
// 		ValidArgs: []string{"bash", "zsh"},
// 		Args:      cobra.OnlyValidArgs,
// 		PreRunE: func(cmd *cobra.Command, args []string) error {
// 			return completionFile.Check(cmd)
// 		},
// 		RunE: func(cmd *cobra.Command, args []string) (err error) {
// 			// generate bash completions, unless zsh is given
// 			if len(args) > 0 && args[0] == "zsh" {
// 				err = rootcmd.GenZshCompletion(completionFile.File)
// 			} else {
// 				err = rootcmd.GenBashCompletion(completionFile.File)
// 			}
// 			outfile.Close()
// 			return
// 		},
// 	}
// 	this := completion
// 	cmd.AddCommand(this)
// 	completionFile = AddFileFlag(this, FileFlagOptions{"out", "o", "output completion script to file",
// 		func(name string) (*os.File, error) {
// 			return os.OpenFile(name, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
// 		},
// 	})
// }
