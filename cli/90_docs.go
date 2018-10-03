// +build !nodocs

package cli

import (
	cf "github.com/ansemjo/aenker/cli/cobraflags"
)

func init() {
	cf.AddGeneratorCommand(RootCommand)
}
