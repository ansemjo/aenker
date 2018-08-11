package cli

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

var (
	infileFlag  string
	infile      io.ReadCloser
	outfileFlag string
	outfile     io.WriteCloser
)

// add optional input/output file flags
func addInFileFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&infileFlag, "input", "i", "", "file to read from instead of stdin")
}
func addOutFileFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&outfileFlag, "output", "o", "", "file to write to instead of stdout")
}

// open input file for reading .. run this in PreRunE
func checkInFileFlag(cmd *cobra.Command, args []string) error {
	if cmd.Flag("input").Changed && infileFlag != "-" {
		f, err := os.Open(infileFlag)
		if err != nil {
			return err
		}
		infile = f
	} else {
		infile = os.Stdin
	}
	return nil
}

// open output file for writing .. run this in PreRunE
func checkOutFileFlag(cmd *cobra.Command, args []string) error {
	if cmd.Flag("output").Changed && outfileFlag != "-" {
		f, err := os.Create(outfileFlag)
		if err != nil {
			return err
		}
		outfile = f
	} else {
		outfile = os.Stdout
	}
	return nil
}
