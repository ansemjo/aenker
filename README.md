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

The format is further described in [SPECIFICATION.md](SPECIFICATION.md).

## example usage

Generate a new random key:

    $ aenker kg -o ~/.aenker

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
    aenker version 0.3.5 (not built with build.go)

Or clone the repository and use the included `build.go` program, which compiles a static binary and
adds more specific information about the built version into the file (`vgo` is required for module
vendoring before `go 1.11.0` is released):

    $ go get golang.org/x/vgo
    $ git clone https://github.com/ansemjo/aenker.git
    $ cd aenker/
    $ make
    vgo mod vendor
    go run build.go -o aenker --tempdir /tmp/aenker-build-tmpgopath
    ./aenker --version
    ./aenker version 0.3.5-4-ge1976f0-dirty
    sha256sum --tag aenker
    SHA256 (aenker) = 0ddd523be3aec435bb044946e1473c651d694658682fabc0b5c665210301b931

If you have [upx](https://upx.github.io/) installed, you can further compress the binary, which
results in approximately a third of the file size:

    $ make compress

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

or

    $ aenker gen completion | sudo tee /usr/share/bash-completion/completions/aenker

## disclaimer

    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
    SOFTWARE.

Please be advised that I am not a professional cryptographer and I made no attempts to produce
constant-time implementations, thus you should probably not use this code for any interactive or
on-the-wire protocols.

This is merely a hobby of mine which I hope can be useful to you.