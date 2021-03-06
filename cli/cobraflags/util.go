// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

// Package cobraflags implements some Flag and Command addons. They are split into
// a seperate package to be able to move them to a sperate repository eventually.
package cobraflags

import (
	"errors"

	"github.com/spf13/cobra"
)

// CheckAll runs pre-run checks of all given checkers.
func CheckAll(cmd *cobra.Command, args []string, checker ...func(*cobra.Command, []string) error) (err error) {
	for _, ch := range checker {
		err = ch(cmd, args)
		if err != nil {
			return
		}
	}
	return
}

// similar to cobra.NoArgs but returns more appropriate error string
func NoArgs(cmd *cobra.Command, args []string) (err error) {
	if len(args) > 0 {
		return errors.New("this command does not accept any positional arguments")
	}
	return
}
