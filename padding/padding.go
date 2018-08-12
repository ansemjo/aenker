// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package padding

//! ATTENTION
//! The functions herein make absolutely no attempt to
//! run in constant time. This likely opens up unwanted
//! timing side-channels when used in an interactive protocol.
//! I intended this package mainly for 'offline' use and
//! archival purposes, so I deem the risk acceptable here.

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
func AddPadding(slice *[]byte, final bool) {

	//! we'll assume that the capacity of the passed slice is the
	//! desired chunksize and reuse that memory
	capacity := cap(*slice)
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
// It returns the length of data as `datalength`, so the caller
// can truncate its slice accordingly.
// The last byte indicates the padding to be expected and
// wether it was a final chunk of a sequence. This information
// is returned as `final`.
// See AddPadding comment for further specifications.
// TODO: should probably return error instead of panicking
func RemovePadding(chunk *[]byte) (final bool) {

	length := len(*chunk)        // get length of chunk
	marker := (*chunk)[length-1] // get last byte, indicating the type
	if marker >= byte(invalid) { // any byte larger than `invalid` is unexpected
		panic("unknown padding type")
	}

	final = marker != byte(running) // was this a final chunk?
	*chunk = (*chunk)[:length-1]    // truncate last byte
	length--

	if marker == byte(padded) { // if type indicates padding ...

		pad := (*chunk)[length-1]          // get the byte that was used to pad
		if !(pad == 0x00 || pad == 0x01) { // any padding byte besides 0 or 1 is unexpected
			panic("unexpected padding byte")
		}

		var offset int          // offset where data ends
		for i := range *chunk { // iterate over bytes, beginning at the end
			if (*chunk)[length-(i+1)] != pad { // and find the first non-pad byte
				offset = i // save the offset
				break
			}
		}
		datalength := length - offset // calulate original slice length
		*chunk = (*chunk)[:datalength]

	}
	return

}
