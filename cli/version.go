package cli

import "strings"

// Version information
// The Format: string will be replaced when downloaded as a git archive.
const version = "0.4"
const commit = "$Format:%h$"

// return the most specific version available
func trueversion() string {
	if strings.Contains(commit, "Format:") {
		return version + " (development)"
	}
	return version + " (commit " + commit + ")"
}
