#!/usr/bin/env python3

# tiny demo decryption script for aenker files
# compatible with files created with aenker 0.5+

import sys
import base64
import argparse

# import cryptography
# $ pip install cryptography
from cryptography.hazmat.backends import default_backend
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.asymmetric import x25519
from cryptography.hazmat.primitives.kdf import hkdf
from cryptography.hazmat.primitives.ciphers import aead


# hkdf info for aenker key derivation
keyinfo = b"aenker elliptic"

# anonymous diffie-hellman and key derivation with hkdf-blake2b512
def keyderivation(private, peer, salt, info):
    private = x25519.X25519PrivateKey._from_private_bytes(private)
    public = x25519.X25519PublicKey.from_public_bytes(peer)
    shared = private.exchange(public)
    return hkdf.HKDF(hashes.BLAKE2b(64), 32, salt, info, default_backend()).derive(shared)


# read header, check magic and return complete header, salt and peer key
def readheader(reader):
    header = reader.read(48)
    if header[:8] != b"aenker\xe7\x9e":
        raise ValueError("unknown magic bytes!")
    salt = header[8:16]
    ephemeral = header[16:]
    return header, salt, ephemeral


# class to decrypt aenker chunkstream
class chunkreader:
    def __init__(self, reader, key, info, chunksize=1984):
        self.reader = reader
        self.cipher = aead.ChaCha20Poly1305(key)
        self.info = info
        self.chunksize = chunksize + 16
        self.i = 0

    # nonce counter, packed in 12 bytes, little-endian
    def ctr(self):
        b = self.i.to_bytes(12, byteorder="little")
        self.i += 1
        return b

    # read next chunk and decrypt
    def next(self):
        chunk = self.reader.read(self.chunksize)
        nonce = self.ctr()
        plain = self.cipher.decrypt(nonce, chunk, self.info)
        return self.removepadding(plain)

    # remove padding from plaintext
    def removepadding(self, chunk):
        typ, chunk = chunk[-1:], chunk[:-1]
        if typ == b"\x00":
            return chunk, False  # running chunk
        elif typ == b"\x01":
            return chunk, True  # final chunk, no padding
        elif typ == b"\x02":
            pad = chunk[-1]  # final chunk, with padding
            return chunk.rstrip(bytes([pad])), True
        raise ValueError("unknown padding type: " + typ)

    # decrypt and write to writer until there is no more
    def decrypt(self, writer):
        final = False
        while not final:
            pt, final = self.next()
            writer.write(pt)


# using private key, decrypt from reader and output to writer
def decrypt_aenker(key, reader, writer):
    private = base64.b64decode(key)
    header, salt, ephemeral = readheader(reader)
    key = keyderivation(private, ephemeral, salt, keyinfo)
    chunkreader(reader, key, header).decrypt(writer)


def cli():

    args = argparse.ArgumentParser()
    args.add_argument("key", help="base64-encoded private key")
    args.add_argument("-i", help="input file", type=argparse.FileType("rb"), default=sys.stdin.buffer)
    args.add_argument("-o", help="output file", type=argparse.FileType("wb"), default=sys.stdout.buffer)
    args = args.parse_args()

    decrypt_aenker(args.key, args.i, args.o)


if __name__ == "__main__":
    cli()
