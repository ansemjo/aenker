# Changelog

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
