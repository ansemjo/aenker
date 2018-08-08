package main

const (
	RunningChunk  = 0
	PaddedChunk   = 1
	UnPaddedChunk = 2
)

func Pad(chunk []byte, size int, last bool) []byte {

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
	return append(chunk, pad...)

}
