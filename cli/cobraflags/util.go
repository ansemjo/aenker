// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

// Package cobraflags implements some Flag and Command addons. They are split into
// a seperate package to be able to move them to a sperate repository eventually.
package cobraflags

import (
	"os"

	"github.com/spf13/cobra"
)

// run pre-run checks of cobra flags
func CheckAll(cmd *cobra.Command, args []string, checker ...func(*cobra.Command, []string) error) (err error) {
	for _, ch := range checker {
		err = ch(cmd, args)
		if err != nil {
			return
		}
	}
	return
}

// Exclusive is a fileopener for FileFlag, which attempts to open the file for
// writing exclusively. I.e. it fails if the file already exists.
func Exclusive(mode os.FileMode) func(name string) (*os.File, error) {
	return func(name string) (*os.File, error) {
		return os.OpenFile(name, os.O_CREATE|os.O_EXCL|os.O_WRONLY, mode)
	}
}
