// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package cobraflags

import (
	"os"

	"github.com/spf13/cobra"
)

type FileFlag struct {
	File *os.File
	Open func(cmd *cobra.Command, args []string) error
}

func AddFileFlag(cmd *cobra.Command, flag, short, usage string,
	open func(string) (*os.File, error), fallback *os.File) (ff *FileFlag) {

	// add flag to command
	str := cmd.Flags().StringP(flag, short, "", usage)

	// return struct with open command for PreRunE
	return &FileFlag{
		Open: func(cmd *cobra.Command, args []string) (err error) {
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

// Exclusive is a fileopener for FileFlag, which attempts to open the file for
// writing exclusively. I.e. it fails if the file already exists.
func Exclusive(mode os.FileMode) func(name string) (*os.File, error) {
	return func(name string) (*os.File, error) {
		return os.OpenFile(name, os.O_CREATE|os.O_EXCL|os.O_WRONLY, mode)
	}
}

// Truncate is a fileopener for FileFlag, which truncates any exiting file or
// creates a new one if it does not exist.
func Truncate(mode os.FileMode) func(name string) (*os.File, error) {
	return func(name string) (*os.File, error) {
		return os.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)
	}
}

// Readonly is a fileopener for FileFlag, which opens an existing file readonly.
func Readonly() func(name string) (*os.File, error) {
	return func(name string) (*os.File, error) {
		return os.Open(name)
	}
}
