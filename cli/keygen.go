// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package cli

import (
	"encoding/base64"

	"github.com/ansemjo/aenker/ae"
	"github.com/spf13/cobra"
)

func init() {
	keygenCmd.Flags().SortFlags = false
	addOutFileFlag(keygenCmd)
}

var keygenCmd = &cobra.Command{
	Use:     "kg",
	Aliases: []string{"keygen"},
	Short:   "Generate a new key",
	Long:    "Generate a random 32-byte key and output base64-encoded form to stdout",
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		// open output file
		// TODO: use the file flag functions from curvekey to enable setting permissions
		err = checkOutFileFlag(cmd, args)
		if err != nil {
			return
		}

		// generate a new key and write encoded form to file
		enc := base64.NewEncoder(base64.StdEncoding, outfile)
		_, err = enc.Write(ae.NewKey())
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
