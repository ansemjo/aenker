// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package cobraflags

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type DirFlag struct {
	Dir   string
	Check func(cmd *cobra.Command, args []string) error
}

func AddDirFlag(cmd *cobra.Command, flag, short, value, usage string) (df *DirFlag) {

	// add flag to command
	str := cmd.Flags().StringP(flag, short, value, usage)

	// return struct with check command for PreRunE
	return &DirFlag{
		Check: func(cmd *cobra.Command, args []string) (err error) {

			stat, err := os.Stat(*str)
			if err != nil {
				return
			}

			if !stat.IsDir() {
				return fmt.Errorf("%s: not a directory", *str)
			}

			df.Dir = *str
			return
		},
	}
}
