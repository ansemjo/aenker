package cli

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	root.AddCommand(flagsDebugCmd)
	addKeyFlags(flagsDebugCmd)
}

var keyfile string
var key []byte

var flagsDebugCmd = &cobra.Command{
	Use:   "flag",
	Short: "flaggy flags",
	Long:  "see how flags work",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return checkKeyFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("key:  % x\n", key)
	},
}

func addKeyFlags(cmd *cobra.Command) {

	cmd.Flags().StringVarP(&keyfile, "keyfile", "f", "", "key file")
	cmd.Flags().BytesBase64VarP(&key, "key", "k", nil, "32 byte key")

}

func checkKeyFlags(cmd *cobra.Command) error {

	kfChg := cmd.Flag("keyfile").Changed
	kyChg := cmd.Flag("key").Changed

	if kfChg && kyChg { // both were given
		return errors.New("only use either one of keyfile or key")
	}

	if !kfChg && !kyChg { // none was given
		return errors.New("one of keyfile or key is required")
	}

	if kyChg && len(key) != 32 { // key was given and it's not 32 bytes
		return errors.New("key must be 32 bytes")
	}

	if kfChg { // keyfile was given

		f, err := os.Open(keyfile) // open keyfile for reading
		if err != nil {
			return err
		}
		defer f.Close()

		line, _, err := bufio.NewReader(f).ReadLine() // read the first line
		if err != nil {
			return err
		}

		n, err := base64.StdEncoding.Decode(line, line) // decode base64 in line
		if err != nil {
			return err
		}

		if n != 32 { // decoded slice needs to be 32 bytes
			return errors.New("key must be 32 bytes")
		}

		key = line[:n] // put it in the key slice
	}

	return nil
}
