package main

import (
	"fmt"
	"os"
)

func stderr(msg interface{}) {
	fmt.Fprintln(os.Stderr, msg)
}

var sfmt = fmt.Sprintf
