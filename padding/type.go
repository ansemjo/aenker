// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package padding

// chunktypetype signals what type of padding
// is to be expected and removed during unpad
type chunktypetype byte

const (
	running  chunktypetype = 0 // a common chunk which is followed by many like itself ...
	unpadded chunktypetype = 1 // a final chunk that fits right inside and needs no padding
	padded   chunktypetype = 2 // a final chunk that requires padding
	invalid  chunktypetype = 3 // any int >= this is not valid
)
