![](anchor.png)

<small>Icon made by [Pixel perfect](https://www.flaticon.com/authors/pixel-perfect) from
[flaticon.com](https://www.flaticon.com/) licensed by
[Creative Commons BY 3.0](http://creativecommons.org/licenses/by/3.0/) </small>

# aenker

`aenker` is a simple commandline utility to encrypt files with an authenticated encryption scheme
([ChaCha20Poly1305](https://cr.yp.to/chacha.html)).

It splits the input into smaller chunks internally and encrypts & authenticates them individually.
Padding and concatenation is done similarly to
[InterMAC](https://rwc.iacr.org/2018/Slides/Hansen.pdf):

```
                 +---------+---------+---+
   chunk + pad + |    1    |    2    | 3 +---+
               | +----+----+-----+---+---+   |
               v      |          |           |
                 +----+----------+-----------v----+   0: running
encrypt + auth + |plainpla|0|nplainpl|0|in.0000002|   1: no padding
               | +-------+------------------------+   2: padding
               v         |           |            |
       +-----------------v-----------v------------v--------------+
       |nonce|mek|z5Q71mOXmu|auth|MXA91ADUiy|auth|5KVhQbzBac|auth|
       +-------^-------------------------------------------------+
               |
               +-+ encrypted random key
```

Authenticated encryption authenticates the ciphertext upon decryption and combined with the above
construction any chunk reordering, bit-flips or even truncation can be detected and are shown as
errors. Only ciphertext that has been successfully decrypted and authenticated is ever written to
the output.

The chunking still alleviates the need to fit the entire file into memory at once or do two passes
over all data. Thus you can also encrypt files of many gigabytes.

## example usage

Generate a new random key:

    $ aenker kg > ~/.aenker

Pack some documents and encrypt them:

    $ tar -czC ~/Documents . | aenker enc -o documents.tar.gz.ae

Decrypt a simple message:

    $ aenker dec -i message.txt.ae
    Hello, World!

Specify a key on the commandline (this is unsafe for various reasons and you should really use
keyfiles generated with `aenker kg`):

    $ aenker e -k lGLDUgFvp8TSwJ17VC9k0/T9mNWvfGoJ42zauMkAFBo= ...

### keyfile

Only the first line of the keyfile is read, so you can add as many comments as you like after that:

    $ cat mykey
    tOnYoytjZpZQjSiEZO0RKYmOZHKJnjmurgKdoJlxB+Y=
    I used this key for x, y and z ...

## installation

You can install from `master` with `go` (make sure `$GOPATH/bin/` is in your `$PATH`):

    $ go get -u github.com/ansemjo/aenker
    $ aenker --version
    aenker version untracked (not built with build.go)

Or clone the repository and use the included `build.go` program, which compiles a static binary and
adds information about the built version into the file (`vgo` is required for module vendoring
before `go 1.11.0` is released):

    $ go get golang.org/x/vgo
    $ git clone https://github.com/ansemjo/aenker.git
    $ cd aenker/
    $ make build
    $ ./aenker --version
    ./aenker version 0.3.2 (0.3.2-4-g102a758-dirty)

Then install in the default user prefix `~/.local/` (make sure `~/.local/bin/` is in your `$PATH`):

    $ make install

Or globally with:

    $ sudo make install PREFIX=/usr/local

### documentation

You can go to [docs/aenker.md](docs/aenker.md) if you're looking at this online.

Installation via the makefile should also install manpages, so `man aenker` should work. Otherwise
use:

    $ aenker gen manual -d /tmp
    $ man -M /tmp aenker

[cobra]: https://github.com/spf13/cobra

All the commands have a nice help message powered by [cobra] aswell, so you can just use `--help` at
any point.

### autocompletion

Completion script generation is also powered by [cobra]. It's available for `bash` and `zsh`.

    $ . <(aenker gen completion)
