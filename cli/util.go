// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package cli

import (
	b64 "encoding/base64"
	"fmt"
	"os"
)

var base64 = b64.StdEncoding.EncodeToString

// Treat any non-nil error as a fatal failure,
// print error to stderr and exit with nonzero status.
func fatal(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "FATAL:", err)
		os.Exit(1)
	}
}
