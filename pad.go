package main

const (
	RunningChunk  = 0
	UnPaddedChunk = 1
	PaddedChunk   = 2
	Invalid       = 3
)

func Pad(chunk []byte, size int, last bool) (out []byte) {

	if !last {
		return append(chunk, RunningChunk)
	}

	// current length
	l := len(chunk)

	// panic if chunk is full
	if l >= size {
		panic("chunk is already >= chunksize")
	}

	// how much padding is needed, may be zero
	needed := size - 1 - l

	// is the last byte null?
	null := chunk[l-1] == 0

	// assemble padding
	pad := make([]byte, needed+1)
	for i := range pad {
		if null {
			pad[i] = 1
		} else {
			pad[i] = 0
		}
	}
	if needed == 0 {
		pad[needed] = UnPaddedChunk
	} else {
		pad[needed] = PaddedChunk
	}

	// append to slice
	out = append(chunk, pad...)
	return

}

//! this is definitely not constant time ...
func Unpad(chunk []byte) (out []byte, last bool) {

	// length of chunk
	l := len(chunk)

	// get last byte
	lb := chunk[l-1]
	if lb >= Invalid {
		panic(sfmt("bad padding: %x", lb))
	}

	// truncate
	out = chunk[:l-1]
	l--

	// is this the last chunk
	last = lb != RunningChunk

	// remove padding if present
	if lb == PaddedChunk {

		// second to last byte
		slb := out[l-1]
		if slb != 0x00 && slb != 0x01 {
			panic(sfmt("bad padding: %x", slb))
		}

		// find first non-padding
		var cut int
		for i := range out {
			if slb != out[l-i-1] {
				cut = i
				break
			}
		}
		out = out[:l-cut]

	}

	return

}
