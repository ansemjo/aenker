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

// Placeholder that is maybe overwritten in pbkdf.go init() if pbkdf build tag given
var AddPbkdfCommand = func(c *cobra.Command) *cobra.Command {
	return c
}

// AddKeygenCommand add the key generator and pubkey converter subcommands to a cobra command.
func AddKeygenCommand(parent *cobra.Command) *cobra.Command {

	var keyfile, comment string

	command := &cobra.Command{
		Use:     "keygen",
		Aliases: []string{"kg", "gen"},
		Short:   "generate a new key",
		Long:    "Generate and save a new random Curve25519 keypair.",
		Example: "  aenker kg -f mykey",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			// format returned errors
			defer func() {
				if err != nil {
					err = fmt.Errorf("aenker keygen: %s", err)
					fatal(err)
				}
			}()

			// generate new random key
			seckey := new([32]byte)
			if _, err = io.ReadFull(rand.Reader, seckey[:]); err != nil {
				return
			}

			// write to file and return pubkey
			pubkey, err := writeKey(seckey, keyfile, comment)
			if err != nil {
				return
			}

			// print info to stdout
			fmt.Printf(`New key saved in %q.
Your public key is: %s
Use the following command to encrypt files for this key:

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
	AddPbkdfCommand(command)

	// add to parent
	parent.AddCommand(command)
	return command
}

// writeKey is the internal function of the keygen, that writes a newly generated key
// to a file with some metadata and comments
func writeKey(key *[32]byte, file, comment string) (pubkey string, err error) {

	// ensure directory exists
	if err = os.MkdirAll(path.Dir(file), 0755); err != nil {
		return
	}

	// create regular files exclusively (will error if it exists)
	create, err := func() (flag int, err error) {
		stat, err := os.Stat(file)
		if os.IsNotExist(err) || (err == nil && stat.Mode().IsRegular()) {
			return os.O_CREATE, nil
		}
		return 0, err
	}()

	// open keyfile
	kf, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC|os.O_EXCL|create, 0600)
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

	// calculate public key and encode to base64
	pubkey = base64(keyderivation.Public(key)[:])

	// append pubkey to header
	header += fmt.Sprintf("# your public key: %s\n", pubkey)

	// format comment if present and append
	if comment != "" {
		header += fmt.Sprintf("# comment: %s\n", strconv.Quote(comment))
	}

	// save secret key to file
	if _, err = kf.WriteString(header + base64(key[:]) + "\n"); err != nil {
		return
	}

	// return publickey
	return

}
