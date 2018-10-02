// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package cobraflags

import (
	"os"

	"github.com/spf13/cobra"
)

type FileFlag struct {
	File *os.File
	Open func(cmd *cobra.Command) error
}

func AddFileFlag(cmd *cobra.Command, flag, short, usage string,
	open func(string) (*os.File, error), fallback *os.File) (ff *FileFlag) {

	// add flag to command
	str := cmd.Flags().StringP(flag, short, "", usage)

	// return struct with open command for PreRunE
	return &FileFlag{
		Open: func(cmd *cobra.Command) (err error) {
			if cmd.Flag(flag).Changed {

				// open given file with passed function
				f, err := open(*str)
				if err != nil {
					return err
				}
				ff.File = f

			} else {
				// if falg wasn't given, use fallback
				ff.File = fallback
			}

			return
		},
	}
}
