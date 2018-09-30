// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package cli

import (
	"os"

	"github.com/spf13/cobra"
)

//type CmdCheckFunc func(cmd *cobra.Command) error

type FileFlag struct {
	File  *os.File
	Check CmdCheckFunc
}

type FileFlagOptions struct {
	Flag  string
	Short string
	Usage string
	Open  func(string) (*os.File, error)
}

func AddFileFlag(cmd *cobra.Command, opts FileFlagOptions) *FileFlag {

	// add flag to command
	str := cmd.Flags().StringP(opts.Flag, opts.Short, "", opts.Usage)
	flag := &FileFlag{}

	// build check function
	flag.Check = func(cmd *cobra.Command) (err error) {
		// if flag was given
		if cmd.Flag(opts.Flag).Changed {
			// open with passed function
			f, err := opts.Open(*str)
			if err != nil {
				return err
			}
			flag.File = f
		}

		return
	}

	return flag

}
