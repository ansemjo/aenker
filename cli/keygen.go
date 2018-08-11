package cli

import (
	"encoding/base64"
	"syscall"

	"github.com/ansemjo/aenker/Aenker"
	"github.com/spf13/cobra"
)

func init() {
	keygenCmd.Flags().SortFlags = false
	addOutFileFlag(keygenCmd)
}

var keygenCmd = &cobra.Command{
	Use:     "kg",
	Aliases: []string{"keygen"},
	Short:   "generate a new key",
	Long:    "Generate a random 32-byte key and output base64-encoded form to stdout",
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		// open output file with restricted (-rw-r----) permissions
		syscall.Umask(0027)
		err = checkOutFileFlag(cmd, args)
		if err != nil {
			return
		}

		// generate a new key and write encoded form to file
		enc := base64.NewEncoder(base64.StdEncoding, outfile)
		_, err = enc.Write(aenker.NewKey())
		if err != nil {
			return
		}
		err = enc.Close()
		if err != nil {
			return
		}
		outfile.Write([]byte{'\n'})
		return outfile.Close()

	},
}
