package cobraflags

import "github.com/spf13/cobra"

// run pre-run checks of cobra flags
func checkAll(cmd *cobra.Command, checker ...func(*cobra.Command) error) (err error) {
	for _, ch := range checker {
		err = ch(cmd)
		if err != nil {
			return
		}
	}
	return
}
