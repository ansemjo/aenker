package cli

import (
	"github.com/spf13/cobra"
)

func init() {
	completionCmd.Flags().SortFlags = false
	addOutFileFlag(completionCmd)
}

var completionCmd = &cobra.Command{
	Use:     "completion [bash|zsh]",
	Aliases: []string{"autocomplete"},
	Short:   "Output autocompletion scripts",
	Long:    "Generate autocompletion scripts to be sourced by your shell.",
	Example: `
Directly source in current shell:
  . <(aenker gen completion)
	
Add to global bash-completions:
  aenker gen completion -o /usr/share/bash-completion/completions/aenker`,
	ValidArgs: []string{"bash", "zsh"},
	Args:      cobra.OnlyValidArgs,
	PreRunE:   checkOutFileFlag,
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		// force cobra to use a simple name, otherwise it could
		// contain path seperator characters (e.g. './aenker')
		rootCmd.Use = "aenker"

		if len(args) > 0 && args[0] == "zsh" { // generate bash completions, unless zsh is given
			err = rootCmd.GenZshCompletion(outfile)
		} else {
			err = rootCmd.GenBashCompletion(outfile)
		}

		outfile.Close()
		return
	},
}
