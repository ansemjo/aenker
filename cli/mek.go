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
		fatal(err)

		fmt.Fprintln(os.Stderr, "MEK:", base64.StdEncoding.EncodeToString(mek))
		os.Stdout.Write(blob)

	},
}

var openmek = &cobra.Command{
	Use:     "openmek",
	Short:   "open a MEK",
	PreRunE: checkKeyFlags,
	Run: func(cmd *cobra.Command, args []string) {

		blob, err := ioutil.ReadAll(os.Stdin)
		fatal(err)

		mek, err := aenker.OpenMEK(key, blob)
		fatal(err)

		fmt.Fprintln(os.Stderr, "MEK:", base64.StdEncoding.EncodeToString(mek))

	},
}

func fatal(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
