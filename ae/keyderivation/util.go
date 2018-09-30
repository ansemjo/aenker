package keyderivation

// get a fixed 32 byte array from slice
func get32(slice []byte) (array *[32]byte) {
	array = new([32]byte)
	copy(array[:], slice)
	return
}
