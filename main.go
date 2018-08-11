package main

import (
	"github.com/ansemjo/aenker/cli"
)

// this is set when building with build.go
var version string

func main() {
	if version != "" {
		cli.SetVersion(version)
	}
	cli.Execute()
}
