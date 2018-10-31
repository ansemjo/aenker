// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package cli

import "strings"

// Version information. The Format: string will be replaced when downloaded as a git archive.
// It is assembled in SpecificVersion().
var Version = "0.5"
var Commit = "$Format:%H$"

// SpecificVersion returns the most specific version available: either '$Version (development)'
// or '$Version (commit $Commit)' - depending on whether Commit was replaced by git upon archive
// creation.
func SpecificVersion() string {
	if strings.Contains(Commit, "Format:") {
		return Version + " (development)"
	}
	return Version + " (commit: " + Commit + ")"
}
