// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

// Package padding adds or removes padding to a plaintext slice.
package padding

import (
	"crypto/subtle"
)

// chunktype signals what type of padding
// is to be expected and removed during unpad
type chunktype byte

// The possible final bytes to indicate the type of this chunk.
const (
	Running  chunktype = '\x00' // a common chunk which is followed by many like itself ...
	Unpadded chunktype = '\x01' // a final chunk that fits right inside and needs no padding
	Padded   chunktype = '\x02' // a final chunk that requires padding
)

// Add appends one or more padding bytes at the end of the slice, depending on whether it is
// a running or a final chunk from a given sequence. If padding is needed, the following rules apply:
//  the last data byte is NOT \x00 --> pad with \x00 bytes
//  the last data byte is \x00     --> pad with \x01 bytes
//
// The very last appended byte indicates the type of chunk and wether padding was applied or not.
// See https://rwc.iacr.org/2018/Slides/Hansen.pdf, page 10 for details.
//
//! WARNING: not constant time, might open up side-channels
func Add(slice *[]byte, final bool, capacity int) {
	// TODO: should probably return error instead of panicking

	length := len(*slice)
	free := capacity - length

	if !final { // if we are not a final slice ...

		if free != 1 { // check that there is space for exactly one byte
			panic("must have exactly one byte free")
		}

		*slice = append(*slice, byte(Running)) // append running chunk marker

	} else {

		if !(free >= 1) { // check that there is space for AT LEAST one byte
			panic("must have at least one byte free")
		}

		var pad byte                    // decide which byte to use for padding
		if (*slice)[length-1] == 0x00 { // if the last data byte is 0x00 ...
			pad = 0x01
		} else {
			pad = 0x00
		}

		needed := free   // how many padding bytes are needed (incl. marker)
		for needed > 1 { // fill all but the last byte with padding ...
			*slice = append(*slice, pad)
			needed--
		}

		if free > 1 { // if we had to use at least one padding byte ...
			*slice = append(*slice, byte(Padded)) // mark this chunk as padded
		} else {
			*slice = append(*slice, byte(Unpadded)) // otherwise unpadded
		}

	}
	return

}

// Remove removes padding which was previously added with Add(). The last byte indicates
// the padding to be expected and wether it was a final chunk of a sequence. This information
// is returned as `final`. See Add() comment for further specifications.
//
//! WARNING: This is my best-effor attempt of creating a constant-time function. Tests with
// https://github.com/oreparaz/dudect do look promising though.
//
// When used with authenticated encryption this might not even be necessary anyway.
func Remove(chunk *[]byte) (final bool) {
	// TODO: should probably return error instead of panicking

	length := len(*chunk)        // get length of chunk
	marker := (*chunk)[length-1] // get last byte, indicating the type
	*chunk = (*chunk)[:length-1] // truncate last byte
	length--

	// final if this was not a 'running' marker
	final = !((subtle.ConstantTimeByteEq(marker, byte(Running)) & 1) == 1)

	// early exit if this is not a final chunk
	// this is not constant time, be we don't want to waste _too_ much time
	// by processing _every_ chunk this way ...
	if !final {
		return
	}

	// mask during pad checking, padding ? 1 : 0
	// when check is set to zero (from the beginning or during the
	// for loop) that means that no future bytes will increment the
	// remove counter, i.e. we're done counting padding bytes.
	check := subtle.ConstantTimeSelect(int(marker&1), 0, 1)

	pad := (*chunk)[length-1] // get the byte that was used to pad
	remove := 0               // number of bytes to be removed

	for i := range *chunk { // iterate over all bytes, beginning at the end
		cur := (*chunk)[length-(i+1)]                              // current byte
		eqok := subtle.ConstantTimeByteEq(cur, pad) & check        // bytes are equal and check is 1
		remove = subtle.ConstantTimeSelect(eqok, remove+1, remove) // increment remove when eqok is 1
		check = eqok                                               // set check to the value of eqok
	}

	data := length - remove  // calulate data length
	*chunk = (*chunk)[:data] // and truncate

	return

}
