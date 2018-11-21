// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package main

import (
	"github.com/ansemjo/aenker/cli"
)

// Version shall be set at compile-time
var Version = "development"

func main() {
	cli.RootCommand.Version = Version
	cli.Execute()
}
