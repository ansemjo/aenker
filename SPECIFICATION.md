# aenker format specification

## Common byte lengths

Some commonly needed byte lengths:

| label    | bytes  | type     | description / where to find                                                                                 |
| -------- | ------ | -------- | ----------------------------------------------------------------------------------------------------------- |
| key      | `32`   | `[]byte` | encryption key,<br> from `golang.org/x/crypto/chacha20poly1305.KeySize`                                     |
| nonce    | `24`   | `[]byte` | extended nonce,<br> from `golang.org/x/crypto/chacha20poly1305.NonceSizeX`                                  |
| overhead | `16`   | `[]byte` | added authentication tag after AEAD encryption, <br>`golang.org/x/crypto/chacha20poly1305.New().Overhead()` |
| plain    | `var`  | `[]byte` | the length of the plaintext                                                                                 |
| chunk    | `8192` | `[]byte` | the length of one (padded) plaintext chunk,<br>is variable and set during encryption                        |
| sealed   | `8208` | `[]byte` | sealed ciphertext chunk, `chunk` + `overhead`                                                               |

## Overview

The data format at rest broadly looks like this:

| offset              | bytes    | type          | description                     |
| ------------------- | -------- | ------------- | ------------------------------- |
| `0`                 | `76`     | `MEKBlob`     | Nonce and sealed media key blob |
| `76 + (i * sealed)` | `sealed` | `SealedChunk` | Sealed ciphertext chunk `i`     |

The total encrypted length calculates as:

    76 + (plain / (chunk - 1)) * (chunk + 16)

## Media Key Blob `MEKBlob`

Two different keys are used in aenker. There is a "key encryption key" `KEK` and a "media encryption
key" `MEK`. The KEK is used to seal the MEK and is provided by the user. The MEK is generated
randomly with every invocation. This is done so that a new unique key is used to encrypt the actual
plaintext every time. Furthermore, the extended version of ChaCha20 is used during key encryption,
so that even if the user provides the same key over and over again, random nonces can safely be
used - the 12 bytes of the regular version would be too short in this case.

The media key blob includes an encrypted media key, an extended nonce for the `XChaCha` cipher, the
encoded chunksize chosen during encryption and the authentication tag.

| offset | bytes | type     | description     |
| ------ | ----- | -------- | --------------- |
| `0`    | `24`  | `[]byte` | XChaCha Nonce   |
| `24`   | `52`  | `[]byte` | Sealed MEK Blob |

The MEK blob is unsealed in [`openMEK()`]. The string `Aenker Media Encryption Key` is used as
associated data during AEAD operations on the MEK blob.

[`openmek()`]: Aenker/mediakey.go#L62

### Unsealed MEK Blob

| offset | bytes | type     | description                                             |
| ------ | ----- | -------- | ------------------------------------------------------- |
| `0`    | `32`  | `[]byte` | Media Encryption Key                                    |
| `32`   | `4`   | `int`    | Chunksize, typecast `uint32`, [encoded] in LittleEndian |

[encoded]: Aenker/util.go#L52

## Chunks

The incoming plaintext is split into equal parts of length `chunk`. To be more precise, it is split
into equal parts of length `chunk - 1` and is then [padded] to the full `chunk` size.

[padded]: padding/padding.go#L34

### Plaintext padding

| offset                | bytes           | type                 | description                    |
| --------------------- | --------------- | -------------------- | ------------------------------ |
| `0`                   | `var`/`plain`   | `[]byte`             | plaintext data                 |
| `chunk - padding - 1` | `var`/`padding` | `0x00`/`0x01`        | padding                        |
| `chunk - 1`           | `1`             | `0x00`/`0x01`/`0x02` | chunk type / padding indicator |

_Assume for this explanation a `chunk` size of `8`._

Given the plaintext bytes:

    67 e6 29 07 2e 2a af fc 5f aa 1e 97 4d aa d3 5d

They are split into three equal parts of length `chunk - 1 = 7`:

    67 e6 29 07 2e 2a af
    fc 5f aa 1e 97 4d aa
    d3 5d

For a "running" chunk, i.e. one that is not at the very end of a sequence, a padding byte `0x00` is
added. For a chunk that needs padding, `0x02` is appended at the _very end_.

| chunk type        | last padding byte |
| ----------------- | ----------------- |
| running           | `0x00`            |
| final, not padded | `0x01`            |
| final, padded     | `0x02`            |

So we get:

    67 e6 29 07 2e 2a af 00
    fc 5f aa 1e 97 4d aa 00
    d3 5d ·· ·· ·· ·· ·· 02

The rest of the bytes in the last chunk are filled with `0x00` if the last data byte is **NOT**
`0x00`, otherwise they are filled with `0x01`:

    67 e6 29 07 2e 2a af 00
    fc 5f aa 1e 97 4d aa 00
    d3 5d 00 00 00 00 00 02

A different example:

    a7 90 c1 1c 41 34 84 41 2d 0b de ca 00

Becomes:

    a7 90 c1 1c 41 34 84 00
    41 2d 0b de ca 00 01 02

#### A note on `chunk` size

The more chunks you have, the more overhead you need to add through the authentication tags. You
might want to raise the chunksize for larger files but keep in mind, that an entire chunk needs to
fit into memory at once. The CLI thus limits the `chunk` size to 1 Gigabyte.

If on the other hand you want to encrypt very small messages, know the message size in advance and
do not care about hiding the size from an adversary you could set the `chunk` size to be the message
size + one byte. That will result in a single chunk.

### Chunk Encryption

During chunk encryption, the regular ChaCha20Poly1305 variant is [used] with the media encryption
key and an incrementing counter as the nonce. The string `Aenker Chunk` appended with the LE-encoded
chunksize is used as [associated data]: i.e. `Aenker Chunk8192` with the default `chunk` size.

[used]: Aenker/chunkstream.go#L44
[associated data]: Aenker/chunkstream.go#L42

#### Sealed Chunks

Naturally, the sealed chunks are longer then the `chunk` size by exactly the amount of AEAD
`overhead`.

| offset  | bytes      | type     | description        |
| ------- | ---------- | -------- | ------------------ |
| `0`     | `chunk`    | `[]byte` | encrypted chunk    |
| `chunk` | `overhead` | `[]byte` | authentication tag |

#### Nonce Counter

The implementation can be seen in [noncecounter.go](Aenker/noncecounter.go). Basically, it is a
simple 64-bit incrementing counter which is encoded in LittleEndian into a 12 byte buffer:

    00 00 00 00 00 00 00 00 00 00 00 00
    01 00 00 00 00 00 00 00 00 00 00 00
    02 00 00 00 00 00 00 00 00 00 00 00
    03 00 00 00 00 00 00 00 00 00 00 00
    ...

The random media encryption keys ensure that the probability of a key-nonce reuse becomes
negligible.
