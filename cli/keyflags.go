package cli

import (
	"bufio"
	"encoding/base64"
	"errors"
	"os"

	"github.com/spf13/cobra"
)

var keyfile string
var key []byte

// add necessary key flags to a command
func addKeyFlags(cmd *cobra.Command) {

	cmd.Flags().StringVarP(&keyfile, "keyfile", "f", "", "file with the key on the first line")
	cmd.Flags().BytesBase64VarP(&key, "key", "k", nil, "encoded key as string")

}

// check and load keys .. run this in PreRunE
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
