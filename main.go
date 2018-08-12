// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

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
