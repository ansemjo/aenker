// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

// +build !nodocs

package cli

import (
	cf "github.com/ansemjo/aenker/cli/cobraflags"
)

func init() {
	cf.AddGeneratorCommand(RootCommand)
}
