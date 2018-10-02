// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package cobraflags

import (
	"os"

	"github.com/spf13/cobra"
)

// some common options when opening files
var (
	CreateExclusive = os.O_CREATE | os.O_EXCL | os.O_WRONLY
)

type FileFlag struct {
	File *os.File
	Open func(cmd *cobra.Command) error
}

type FileFlagOptions struct {
	Flag    string
	Short   string
	Usage   string
	Open    func(string) (*os.File, error)
	Default *os.File
}

func AddFileFlag(cmd *cobra.Command, opts FileFlagOptions) *FileFlag {

	// add flag to command
	str := cmd.Flags().StringP(opts.Flag, opts.Short, "", opts.Usage)
	flag := &FileFlag{}

	// build check function
	flag.Open = func(cmd *cobra.Command) (err error) {
		// if flag was given
		if cmd.Flag(opts.Flag).Changed {
			// open with passed function
			f, err := opts.Open(*str)
			if err != nil {
				return err
			}
			flag.File = f
		} else {
			flag.File = opts.Default
		}
		return
	}

	return flag

}
