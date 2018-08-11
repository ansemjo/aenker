package cli

import (
	"github.com/ansemjo/aenker/Aenker"
	"github.com/spf13/cobra"
)

func init() {
	encryptCmd.Flags().SortFlags = false
	addKeyFlags(encryptCmd)
	addInFileFlag(encryptCmd)
	addOutFileFlag(encryptCmd)
	addChunkSizeFlag(encryptCmd)
}

var encryptCmd = &cobra.Command{
	Use:     "enc",
	Aliases: []string{"e", "encrypt"},
	Short:   "encrypt a file",
	Long:    "encrypt stdin and place the ciphertext in stdout",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {

		err = parseChunkSize(cmd, args)
		if err != nil {
			return
		}

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
		_, err = ae.Encrypt(outfile, infile, chunksize)
		infile.Close()
		outfile.Close()
		fatal(err)
		return

	},
}
