package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/ansemjo/aenker/aenker"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(decryptCmd)
	decryptCmd.Flags().SortFlags = false
	addKeyFlags(decryptCmd)
	addChunkSizeFlag(decryptCmd)
}

var decryptCmd = &cobra.Command{
	Use:   "d",
	Short: "decrypt a file",
	Long:  "decrypt stdin and place the plaintext in stdout",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {

		err = parseChunkSize(cmd, args)
		if err != nil {
			return
		}

		return checkKeyFlags(cmd, args)

	},
	Run: func(cmd *cobra.Command, args []string) {

		reader := os.Stdin
		writer := os.Stdout

		// TODO: handle meks transparently in aenker.*crypt
		blob := make([]byte, aenker.MekBlobSize)
		_, err := io.ReadFull(reader, blob)
		fatal(err)

		mek, err := aenker.OpenMEK(key, blob)
		fatal(err)

		ae, err := aenker.NewAenker(mek, chunksize)
		fatal(err)

		lw, err := ae.Decrypt(writer, reader)
		fatal(err)

		fmt.Fprintf(os.Stderr, "wrote %d bytes\n", lw)

	},
}
