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
		Short:   "generate a new keypair",
		Long:    "Generate and save a new random Curve25519 keypair.",
		Example: "  aenker kg -p publickey -o secretkey",
		// PreRunE: func(cmd *cobra.Command, args []string) (err error) {

		// 	// output file
		// 	err = private.Open(cmd, args)
		// 	if err != nil {
		// 		return
		// 	}

		// 	// public key file
		// 	err = public.Open(cmd, args)
		// 	if err != nil {
		// 		os.Remove(private.File.Name())
		// 		return
		// 	}

		// 	return
		// },
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

			// open pubkeyfile
			pf, err := os.OpenFile(keyfile+".pub", os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				return
			}
			defer pf.Close()

			// generate some metadata
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
			metadata := fmt.Sprintf("# aenker key: %s@%s, %s\n", username, hostname, timestamp)

			// maybe quote and append the comment
			if cmd.Flag("comment").Changed {
				metadata += fmt.Sprintf("# comment: %s\n", strconv.Quote(comment))
			}

			// generate new random key
			key := new([32]byte)
			_, err = io.ReadFull(rand.Reader, key[:])
			fatal(err)

			// calculate public key
			pub := keyderivation.Public(key)

			// replace with base64
			keystr, pubstr := base64(key[:]), base64(pub[:])

			// save to files
			if _, err = kf.WriteString(metadata + keystr + "\n"); err != nil {
				return
			}
			if _, err = pf.WriteString(metadata + pubstr + "\n"); err != nil {
				return
			}

			// print info to stdout
			fmt.Printf(`new keypair saved in %q
your public key is: %s
use the following command to encrypt files for this key:

  aenker seal -p %s ...

`, keyfile+"{,.pub}", pubstr, pubstr)

			return
		},
	}
	command.Flags().SortFlags = false

	// try to assemble default keyfile path
	var defaultkey string
	if home, err := os.UserHomeDir(); err != nil {
		defaultkey = path.Join("./", "aenker") // fallback to current dir
	} else {
		defaultkey = path.Join(home, ".local", "share", "aenker", "aenker")
	}

	// define flags for parsing
	command.Flags().StringVarP(&keyfile, "file", "f", defaultkey, "save key to this file")
	command.Flags().StringVarP(&comment, "comment", "c", "", "add comment to keyfile")

	// add subcommands
	AddPubkeyCommand(command)
	if AddPbkdfCommand != nil {
		AddPbkdfCommand(command)
	}
	parent.AddCommand(command)
	return command
}
