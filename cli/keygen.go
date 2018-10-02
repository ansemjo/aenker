// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

// +build !nokeygen

package cli

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"

	cf "github.com/ansemjo/aenker/cli/cobraflags"
	"github.com/ansemjo/aenker/keyderivation"
	"github.com/spf13/cobra"
)

func init() {
	AddKeygenCommand(RootCommand)
}

func AddKeygenCommand(parent *cobra.Command) {

	var private *cf.FileFlag
	var public *cf.FileFlag

	command := &cobra.Command{
		Use:     "keygen",
		Aliases: []string{"kg"},
		Short:   "generate a new keypair",
		Long:    "Generate and save a new Curve25519 keypair.",
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {

			// output file
			err = private.Open(cmd)
			if err != nil {
				return
			}

			// public key file
			err = public.Open(cmd)
			if err != nil {
				os.Remove(private.File.Name())
				return
			}

			return
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			// close all open file upon exit
			defer func() {
				for _, f := range []*cf.FileFlag{private, public} {
					if f.File != nil {
						f.File.Close()
					}
				}
			}()

			// generate a new key
			key := new([32]byte)
			_, err = io.ReadFull(rand.Reader, key[:])
			fatal(err)

			// write encoded key to file
			if private.File == os.Stdout {
				fmt.Fprintf(os.Stderr, "private key:\n  ")
			}
			_, err = fmt.Fprintln(private.File, base64(key[:]))
			fatal(err)

			// if public was given, write public part
			if public.File != nil {

				pub := keyderivation.Public(key)
				if public.File == os.Stdout {
					fmt.Fprintf(os.Stderr, "public key:\n  ")
				}
				_, err = fmt.Fprintln(public.File, base64(pub[:]))

			}

			return
		},
	}
	command.Flags().SortFlags = false

	// file opener
	exclusive := func(mode os.FileMode) func(name string) (*os.File, error) {
		return func(name string) (*os.File, error) {
			return os.OpenFile(name, os.O_CREATE|os.O_EXCL|os.O_WRONLY, mode)
		}
	}

	// add the output file flags
	private = cf.AddFileFlag(command, "out", "o", "write output to file (default: stdout)", exclusive(0600), os.Stdout)
	public = cf.AddFileFlag(command, "pubkey", "p", "write public key to file (default: stdout)", exclusive(0644), os.Stdout)

	AddPubkeyCommand(command)
	parent.AddCommand(command)
}
