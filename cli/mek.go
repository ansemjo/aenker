package cli

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ansemjo/aenker/aenker"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(sealmek)
	rootCmd.AddCommand(openmek)
	sealmek.Flags().SortFlags = false
	openmek.Flags().SortFlags = false
	addKeyFlags(sealmek)
	addKeyFlags(openmek)

}

var sealmek = &cobra.Command{
	Use:     "sealmek",
	Short:   "seal a new MEK",
	PreRunE: checkKeyFlags,
	Run: func(cmd *cobra.Command, args []string) {

		mek, blob, err := aenker.NewMEK(key)
		if err == nil {
			fmt.Fprintln(os.Stderr, "MEK:", base64.StdEncoding.EncodeToString(mek))
			os.Stdout.Write(blob)
		}

	},
}

var openmek = &cobra.Command{
	Use:     "openmek",
	Short:   "open a MEK",
	PreRunE: checkKeyFlags,
	Run: func(cmd *cobra.Command, args []string) {

		blob, err := ioutil.ReadAll(os.Stdin)
		if err == nil {
			mek, err := aenker.OpenMEK(key, blob)
			if err == nil {
				fmt.Fprintln(os.Stderr, "MEK:", base64.StdEncoding.EncodeToString(mek))
			}
		}

	},
}
