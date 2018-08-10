package cli

import (
	"fmt"

	"github.com/c2h5oh/datasize"
	"github.com/spf13/cobra"
)

var chunksize datasize.ByteSize
var chunksizeArg string

// add necessary key flags to a command
func addChunkSize(cmd *cobra.Command) {

	cmd.Flags().StringVar(&chunksizeArg, "chunksize", "", "size of padded plaintext chunk")

}

func parseChunkSize(cmd *cobra.Command) (err error) {

	if cmd.Flag("chunksize").Changed {

		err = chunksize.UnmarshalText([]byte(chunksizeArg))
		if err == nil {
			fmt.Println(chunksize)
		}

	}
	return

}
