package cli

import (
	"github.com/ansemjo/aenker/Aenker"
	"github.com/spf13/cobra"
)

func init() {
	decryptCmd.Flags().SortFlags = false
	addKeyFlags(decryptCmd)
	addInFileFlag(decryptCmd)
	addOutFileFlag(decryptCmd)
}

var decryptCmd = &cobra.Command{
	Use:   "d",
	Short: "decrypt a file",
	Long:  "decrypt stdin and place the plaintext in stdout",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {

		err = checkKeyFlags(cmd, args)
		if err != nil {
			return
		}

		err = checkInFileFlag(cmd, args)
		if err != nil {
			return
		}

		err = checkOutFileFlag(cmd, args)
		if err != nil {
			return
		}

		return
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		ae := aenker.NewAenker(&key)
		_, err = ae.Decrypt(outfile, infile)
		infile.Close()
		outfile.Close()
		fatal(err)
		return

	},
}
