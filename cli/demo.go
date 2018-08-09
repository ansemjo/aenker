package cli

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ansemjo/aenker/aenker"
	"golang.org/x/crypto/chacha20poly1305"
)

func init() {
	root.AddCommand(demo)
}

var demo = &cobra.Command{
	Use: "demo",
	Run: func(cmd *cobra.Command, args []string) {
		oldmain()
	},
}

func fatal(e error) {
	if e != nil {
		panic(e)
	}
}

func oldmain() {

	reader := strings.NewReader("Clear is better than clever.\n")
	// stderr(sfmt("reader has %d bytes", reader.Len()))

	writer := os.Stdout
	var midbuf bytes.Buffer
	var outbuf bytes.Buffer

	zerokey := make([]byte, chacha20poly1305.KeySize)
	ae, err := aenker.NewAenker(zerokey)
	fatal(err)

	_, err = ae.Encrypt(&midbuf, reader)
	fatal(err)

	// add extra data
	midbuf.WriteString("blablabla")

	_, err = ae.Decrypt(&outbuf, &midbuf)
	if err == aenker.ErrExtraData {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("WARN: %s", err.Error()))
	} else {
		fatal(err)
	}

	io.Copy(writer, &outbuf)

}
