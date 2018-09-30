// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package padding

import (
	"crypto/subtle"
)

// chunktypetype signals what type of padding
// is to be expected and removed during unpad
type chunktypetype byte

const (
	running  chunktypetype = 0 // a common chunk which is followed by many like itself ...
	unpadded chunktypetype = 1 // a final chunk that fits right inside and needs no padding
	padded   chunktypetype = 2 // a final chunk that requires padding
)

// AddPadding appends one or more padding bytes at the end of
// the slice, depending on whether it is a running or a
// final chunk from a given sequence. If padding is needed, the
// following rules apply:
//
// - the last data byte is NOT 0x00 --> pad with 0x00 bytes
// - the last data byte is 0x00 --> pad with 0x01 bytes
//
// The very last appended byte indicates the type of chunk and
// wether padding was applied or not.
//
// See https://rwc.iacr.org/2018/Slides/Hansen.pdf, page 10.
// TODO: should probably return error instead of panicking
//! WARN: not constant time, might open up side-channels
func AddPadding(slice *[]byte, final bool, capacity int) {

	//! we'll assume that the capacity of the passed slice is the
	//! desired chunksize and reuse that memory
	length := len(*slice)
	free := capacity - length

	if !final { // if we are not a final slice ...

		if free != 1 { // check that there is space for exactly one byte
			panic("must have exactly one byte free")
		}

		*slice = append(*slice, byte(running)) // append running chunk marker

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
			*slice = append(*slice, byte(padded)) // mark this chunk as padded
		} else {
			*slice = append(*slice, byte(unpadded)) // otherwise unpadded
		}

	}
	return

}

// RemovePadding looks at a chunk to see how much padding must
// be removed, which was previously added with AddPadding.
// The last byte indicates the padding to be expected and
// wether it was a final chunk of a sequence. This information
// is returned as `final`.
// See AddPadding comment for further specifications.
// TODO: should probably return error instead of panicking
//! WARNING: This is my best-effor attempt of creating a constant-
//! time function. Tests with oreparaz/dudect do look promising though.
func RemovePadding(chunk *[]byte) (final bool) {

	length := len(*chunk)        // get length of chunk
	marker := (*chunk)[length-1] // get last byte, indicating the type
	*chunk = (*chunk)[:length-1] // truncate last byte
	length--

	// final if this was not a 'running' marker
	final = !((subtle.ConstantTimeByteEq(marker, byte(running)) & 1) == 1)

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
