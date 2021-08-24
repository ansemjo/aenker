# aenker

[![codebeat badge](https://codebeat.co/badges/0a98d937-6695-4dc1-ba6f-c439226bea01)](https://codebeat.co/projects/github-com-ansemjo-aenker-master)

`aenker` is a simple commandline utility to encrypt files to a public key ([Curve25519][0]) with an
authenticated encryption scheme ([ChaCha20Poly1305][1]). This is basically an [ECIES][2].

The input is split into smaller chunks internally and is encrypted & authenticated individually.
Padding and concatenation is done similarly to [InterMAC][3]. The key used for encryption is derived
with [HKDF][4] using [Blake2b][5] after performing anonymous Diffie-Hellman with a given public and
a random ephemeral private key. All this is further described in the
[specification](SPECIFICATION.md).

[0]: https://cr.yp.to/ecdh.html
[1]: https://tools.ietf.org/html/rfc7539
[2]: https://en.wikipedia.org/wiki/Integrated_Encryption_Scheme
[3]: https://rwc.iacr.org/2018/Slides/Hansen.pdf
[4]: https://tools.ietf.org/html/rfc5869
[5]: https://blake2.net/

![](assets/overview.png)

Authenticated encryption authenticates the ciphertext upon decryption and combined with the above
construction any chunk reordering, bit-flips or even truncation can be detected and are shown as
errors. Only ciphertext that has been successfully decrypted and authenticated is ever written to
the output. The chunking still alleviates the need to fit the entire file into memory at once or do
two passes over all data. Thus you can also encrypt files of many gigabytes.

## INSTALLATION

### Install directly with `go`:

    go get -u github.com/ansemjo/aenker

### Install a binary release / package

Download a release [from GitHub](https://github.com/ansemjo/aenker/releases) (replace 
`$VERSION` with the desired release):

    curl -Lo aenker https://github.com/ansemjo/aenker/releases/download/$VERSION/aenker-linux-amd64
    chmod +x aenker
    ./aenker --help


### Compile from sources:

Download a [tarball from GitHub](https://github.com/ansemjo/aenker/archive/master.tar.gz) and
use the included `makefile` to build a static binary and embed proper version information:

    cd aenker-master/
    make
    make install PREFIX=~/.local


## USAGE

### Key Generation

First, you need a keypair. To generate a new random keypair use the builtin keygenerator:

    aenker keygen [-f where/to/store/seckey]

Without any arguments, this will store the key in the default location `~/.local/share/aenker/aenkerkey`
and your public key will be printed to the terminal. Send your **public** key to anyone who wants to
encrypt data for you and keep your private key .. well, private.

If you want to display your public key later or calculate the public key to a given private key, you
can use the subcommand `show`:

    aenker show [-k path/to/seckey]

**Note:** aenker only performs anonymous Diffie-Hellman and the keys are not signed or certified. To
protect against man-in-the-middle attacks you should transfer the key over a secure channel or verify
the integrity on a different channel.

### Encryption / Decryption

Encrypt a simple message using the public key with the subcommand `seal`:

    echo 'Hello, World!' | aenker seal -p lGLD...AFBo= > message.ae

Decrypt messages with the `open` subcommand. If your key is stored at the default location you can
decrypt a message by simply piping the encrypted message into aenker:

    aenker open [-k path/to/seckey] < message.ae

Input and output files can be specified with the `-i` and `-o` flags respectively. The terms `seal`
and `open` are commonly used in the context of AEADs but you can also use their aliases `encrypt`
and `decrypt` if you prefer:

    aenker decrypt -i documents.tar.ae -k mykey | tar -xf -

The key flags `-p`/`--peer` and `-k`/`--key` accept the base64-encoded keys on the commandline or
the name of a file which contains the key alone on one line. Specifically, the first match to the
regular expression `/^[A-Za-z0-9+/]{43}=$/` is used, so you can add as many comments as you like to
your key files.

Specifying the key on the commandline is convenient for public keys but should be avoided for
private keys:

    ... | aenker seal -p lGLDUgFvp8TSwJ17VC9k0/T9mNWvfGoJ42zauMkAFBo= > message.ae

### Advanced Key Generation

Generally, Curve25519 - and thus aenker - accepts any 32 byte value as a key. You could generate a
private key by any other means and then only calculate the public key to distribute it. Possibilities
include:

* reading 32 bytes of system randomness from `/dev/urandom`
* use an implementation of Argon2i to derive a key from a password, i.e.
  [ansemjo/stdkdf](https://github.com/ansemjo/stdkdf)
* ...

## DOCUMENTATION

All of the commands output a nicely formatted help message, so you can use `help` at any time:

    aenker help

If you prefer, you can instead install and read manpages with:

    aenker docs man -d ~/.local/share/man/
    man aenker

Completion scripts for your shell can be generated and sourced with:

     . <(aenker docs completion)

### File Detection

Append this piece to your `~/.magic` file:

    0 string aenker\xe7\x9e aenker encrypted file
    !:mime application/octet-stream

And `file(1)` should detect encrypted files as `aenker encrypted file`.

## DISCLAIMER

Please be advised that I am not a professional cryptographer. This is merely a hobby of mine which I
hope can be useful to you.

    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
    SOFTWARE.
