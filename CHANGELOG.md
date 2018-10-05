# Changelog

## 0.4.0

Major overhaul. Aenker is now an ECIES implementation, where curve25519 is used to perform anonymous
diffie-hellman and an encryption key is derived using HKDF.

the `build.go` script is ditched and releases are built using
[mkr](https://github.com/ansemjo/makerelease).

Please consult the readme and specification at this tag to find out more. Here is the commit
overview since 0.3.6:

```
f485870 silence directory "changes" during release builds
43180fe update readme and specification
1250e68 add quasi-symmetric encryption mode
ca95ca4 logic fix in chunkreader to handle truncated ciphertext
83741d4 add input/output file flags to encrypt and decrypt
f4dd216 add package and method comments in cli
7880501 add a few missing license headers
4ab5703 fix documentation/autocompletion commands, lexicographical init sort
2b1a29c add proper version information with gitattributes
4ff2034 improve keygen sumcommand and tweak FileFlag
28e66dd fix cli in simplified form for now
70b51e1 add some missing license headers
deeb2b2 [broken] tidy up and document keyderivation and padding packages
8f74f0b [broken] tidy up and document chunkstream package
38b78f2 [broken cli] finish ae package with documentation
f361abe [cli broken] new header format with only elliptic kdf support
e861700 rename packages and remove old chunkstreamer
e1f6378 implement a chunkreader which implements io.Reader
48133a3 write a chunkwriter, which implements a simple WriteCloser interface
aec4568 remove local debug files
a7e5720 fix the chunkstream aead at package level
b532659 define noncecounter output size upon creation
2710beb added vscode debug launch configuration but don't track it
caf5750 refactor keyderivation a little
35fbb6e write the key derivation functions
fd48d4c use the chunkstream subpackage
fbbe564 [broken] add new header decoder and move files to subpackages
a55ace8 rename aenker internals as ae
2e0023f remove syscall import
516103e add OS and ARCH in build target
b254527 remove build.go dependency
38067ac remove docs for now
f765e60 added key derivation func for upcoming v2 format
357e4cc added spreadsheet to calculate and visualize overhead
a228d24 add metadata to mkr prepare and use the same build.go script
8abd584 update mkrelease targets
40b7630 make makefile compatible with ansemjo/makerelease
```

## 0.3.6

### Added specification

I attempted to document the on-disk format in [SPECIFICATION.md](SPECIFICATION.md).

### Reproducible builds

By adding a static `--tempdir ...` argument to the `build.go` invocation, the builds became
reproducible.

## 0.3.5

### `RemovePadding()` was refactored into a constant-time function

Although upon further inspection the `RemovePadding()` function does not seem to open up any
side-channels, as it is only run on already-authenticated plaintext, it still seems to be good form.
Especially if this code will ever be used in a different context.

Constant-time execution is achieved with functions from `crypto/subtle` and logic operations. Most
importantly, the loop does not break early but always iterates over all bytes.

Some testing was performed with [oreparaz/dudect](https://github.com/oreparaz/dudect) (or my forked
[ansemjo/dudect](https://github.com/ansemjo/dudect) respectively, which
[adds support](https://github.com/oreparaz/dudect/pull/3) for timing Go functions). The results look
promising:

```
$ make -B go && ./dudect_go_-O2
clang -O2 -Iinc/ -c src/cpucycles.c -o src/cpucycles.o
clang -O2 -Iinc/ -c src/fixture.c -o src/fixture.o
clang -O2 -Iinc/ -c src/random.c -o src/random.o
clang -O2 -Iinc/ -c src/ttest.c -o src/ttest.o
clang -O2 -Iinc/ -c src/percentile.c -o src/percentile.o
go build -o dut/go/dut_go.so -buildmode=c-shared dut/go/dut_go.go dut/go/random.go dut/go/makeslice.go dut/go/measure.go
clang -O2 -Iinc/ -o dudect_go_-O2 dut/go/dut_go.c src/cpucycles.o src/fixture.o src/random.o src/ttest.o src/percentile.o dut/go/dut_go.so -lm
meas:    0.10 M, max t:   +2.90, max tau: 9.28e-03, (5/tau)^2: 2.90e+05. For the moment, maybe constant time.
meas:    0.20 M, max t:   +1.71, max tau: 3.86e-03, (5/tau)^2: 1.68e+06. For the moment, maybe constant time.
meas:    0.09 M, max t:   +1.48, max tau: 4.91e-03, (5/tau)^2: 1.04e+06. For the moment, maybe constant time.
meas:    0.36 M, max t:   +1.71, max tau: 2.84e-03, (5/tau)^2: 3.11e+06. For the moment, maybe constant time.
meas:    0.39 M, max t:   +1.77, max tau: 2.83e-03, (5/tau)^2: 3.11e+06. For the moment, maybe constant time.
meas:    0.47 M, max t:   +1.43, max tau: 2.09e-03, (5/tau)^2: 5.72e+06. For the moment, maybe constant time.
meas:    0.69 M, max t:   +2.29, max tau: 2.76e-03, (5/tau)^2: 3.27e+06. For the moment, maybe constant time.
meas:    0.78 M, max t:   +1.84, max tau: 2.08e-03, (5/tau)^2: 5.76e+06. For the moment, maybe constant time.
meas:    0.68 M, max t:   +1.50, max tau: 1.82e-03, (5/tau)^2: 7.56e+06. For the moment, maybe constant time.
meas:    0.76 M, max t:   +1.61, max tau: 1.85e-03, (5/tau)^2: 7.33e+06. For the moment, maybe constant time.
meas:    1.02 M, max t:   +1.55, max tau: 1.53e-03, (5/tau)^2: 1.06e+07. For the moment, maybe constant time.
meas:    1.11 M, max t:   +2.01, max tau: 1.90e-03, (5/tau)^2: 6.89e+06. For the moment, maybe constant time.
meas:    0.77 M, max t:   +1.97, max tau: 2.25e-03, (5/tau)^2: 4.94e+06. For the moment, maybe constant time.
meas:    0.83 M, max t:   +2.13, max tau: 2.34e-03, (5/tau)^2: 4.58e+06. For the moment, maybe constant time.
meas:    0.89 M, max t:   +2.40, max tau: 2.55e-03, (5/tau)^2: 3.85e+06. For the moment, maybe constant time.
meas:    0.95 M, max t:   +2.62, max tau: 2.70e-03, (5/tau)^2: 3.44e+06. For the moment, maybe constant time.
[...]
```

As the original author [notes](https://github.com/oreparaz/dudect#typical-output), we're _probably_
safe if the `t` value never goes beyond 10. It is, however, **not** a guarantee. The following was
used to measure execution time:

```golang
// measure.go

package main

import (
	"github.com/ansemjo/aenker/padding"
)

// if you change these values, change them in dut_go.c aswell!
const chunksize = 8 * 1024
const measurements = 1e5

// prepare byte slices with padded data
func prepareData(data *[][]byte, classes *[]byte) {

	for i := range *data {

		(*classes)[i] = randombit()

		var length int
		if (*classes)[i] == 1 {
			length = randomint(1, chunksize-1)
		} else {
			length = chunksize - 1
		}

		slice := (*data)[i][:length]
		randombytes(&slice)

		padding.AddPadding(&slice, true, chunksize)

	}

}

// do work, remove padding
func doWork(data *[]byte) int {

	padding.RemovePadding(data)

	return 0

}
```
