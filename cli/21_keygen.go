// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

// +build !nokeygen

package cli

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"os/user"
	"path"
	"strconv"
	"time"

	"github.com/ansemjo/aenker/keyderivation"
	"github.com/spf13/cobra"
)

func init() {
	AddKeygenCommand(RootCommand)
}

// Placeholder that is maybe declared in pbkdf.go init() if pbkdf build tag given
var AddPbkdfCommand func(*cobra.Command) *cobra.Command

// AddKeygenCommand add the key generator and pubkey converter subcommands to a cobra command.
func AddKeygenCommand(parent *cobra.Command) *cobra.Command {

	var keyfile, comment string

	command := &cobra.Command{
		Use:     "keygen",
		Aliases: []string{"kg", "gen"},
		Short:   "generate a new key",
		Long:    "Generate and save a new random Curve25519 keypair.",
		Example: "  aenker kg -p publickey -o secretkey",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			// format returned errors
			defer func() {
				if err != nil {
					err = fmt.Errorf("aenker keygen: %s", err)
					fatal(err)
				}
			}()

			// ensure directory exists
			if err = os.MkdirAll(path.Dir(keyfile), 0755); err != nil {
				return
			}

			// open keyfile
			kf, err := os.OpenFile(keyfile, os.O_EXCL|os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				return
			}
			defer kf.Close()

			// generate some metadata to embed in keyfile
			username := func() string {
				u, e := user.Current()
				if e == nil {
					return u.Username
				}
				return "unknown"
			}()
			hostname := func() string {
				h, e := os.Hostname()
				if e == nil {
					return h
				}
				return "unknown"
			}()
			timestamp := time.Now().UTC().Format(time.RFC3339)

			// prepare a file header from metadata
			header := fmt.Sprintf("# aenker secret key: %s@%s, %s\n", username, hostname, timestamp)

			// generate new random key
			seckey := new([32]byte)
			_, err = io.ReadFull(rand.Reader, seckey[:])
			fatal(err)

			// calculate public key and encode to base64
			pubkey := base64(keyderivation.Public(seckey)[:])

			// append pubkey to header
			header += fmt.Sprintf("# your public key: %s\n", pubkey)

			// prepare comment if present and append
			if cmd.Flag("comment").Changed {
				header += fmt.Sprintf("# comment: %s\n", strconv.Quote(comment))
			}

			// save secret key to file
			if _, err = kf.WriteString(header + base64(seckey[:]) + "\n"); err != nil {
				return
			}

			// print info to stdout
			fmt.Printf(`new key saved in %q
your public key is: %s
use the following command to encrypt files for this key:

  aenker seal -p %s ...

`, keyfile, pubkey, pubkey)

			return
		},
	}
	command.Flags().SortFlags = false

	// define flags for parsing
	command.Flags().StringVarP(&keyfile, "file", "f", defaultkey, "save key to this file")
	command.Flags().StringVarP(&comment, "comment", "c", "", "add comment to keyfile")

	// add subcommands
	AddPubkeyCommand(command)
	if AddPbkdfCommand != nil {
		AddPbkdfCommand(command)
	}

	// add to parent
	parent.AddCommand(command)
	return command
}
