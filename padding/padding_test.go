package padding

import (
	"bytes"
	"testing"
)

// create a byteslice with the given capacity from a string
func slice(str string, c int) []byte {
	sl := make([]byte, len(str), c)
	copy(sl, str)
	return sl
}

func TestAdd(t *testing.T) {

	// positive test cases that should produce expected results
	positive := []struct {
		from, to string
		cap      int
		final    bool
	}{
		{"Hello, World!", "Hello, World!\x00\x00" + string(Padded), 16, true},
		{"\x00\x00\x00", "\x00\x00\x00" + string(Running), 4, false},
		{"nil\x00", "nil\x00\x01" + string(Padded), 6, true},
		{"nil\x00", "nil\x00" + string(Running), 5, false},
		{"nil\x00", "nil\x00" + string(Unpadded), 5, true},
		{"0", "0\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00" + string(Padded), 64, true},
		{"\x00", "\x00\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01" + string(Padded), 64, true},
	}

	for i, tc := range positive {
		from := slice(tc.from, tc.cap)
		to := slice(tc.to, tc.cap)
		err := Add(&from, tc.final, tc.cap)
		if err != nil {
			t.Errorf("pos[%d] Unexpected error: %s", i, err.Error())
			continue
		}
		ok := bytes.Compare(from, to)
		if ok != 0 {
			t.Errorf("pos[%d] Padding '%x' did not yield the expected result.", i, tc.from)
			continue
		}
	}

	// negative test cases that should fail with an expected error
	negative := []struct {
		slice []byte
		final bool
		err   string
	}{
		{slice("Hello, World!", 13), false, ErrOneByte},
		{slice("Hello, World!", 15), false, ErrOneByte},
		{slice("Hello, World!", 13), true, ErrSize},
	}

	for i, tc := range negative {
		err := Add(&tc.slice, tc.final, cap(tc.slice))
		if err == nil {
			t.Errorf("neg[%d] Unexpected success?!", i)
			continue
		}
		if err.Error() != tc.err {
			t.Errorf("neg[%d] Unexpected error: %s", i, err.Error())
			continue
		}
	}

}
