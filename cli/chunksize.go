package cli

import (
	ce "github.com/ansemjo/aenker/error"
	"github.com/c2h5oh/datasize"
	"github.com/spf13/cobra"
)

var chunksize int
var chunksizeArg string

// add chunksize flag to a command
func addChunkSizeFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(&chunksizeArg, "chunksize", "8kB", "plaintext chunks")
}

// convert datasize string to an int, run in PreRunE
func parseChunkSize(cmd *cobra.Command) (err error) {

	var size datasize.ByteSize
	err = size.UnmarshalText([]byte(chunksizeArg))
	if err != nil {
		return
	}

	if size > datasize.GB {
		return ce.ConstError("chunksize too large")
	}
	chunksize = int(size.Bytes())

	return

}
