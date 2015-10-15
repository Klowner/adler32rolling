package adler32rolling

import (
	"strings"
	"testing"
)

// checksum is a slow but simple implementation of the Adler-32 checksum.
// It is a straight port of the sample code in RFC 1950 section 9 except
// it is initialized with 0 rather than 1
func checksum(p []byte) uint32 {
	s1, s2 := uint32(0), uint32(0)
	for _, x := range p {
		s1 = (s1 + uint32(x)) % mod
		s2 = (s2 + s1) % mod
	}
	return s2<<16 | s1
}

var rollStrings = []string{
	"You're a lizard, Harry.",
	"Don't tell me you've never seen a leprechaun before.",
	strings.Repeat("Too legit, too legit to quit, (hay haaaay).", 1e2),
	strings.Repeat("ABCDEFGHIJKLMNOPQRSTUVWXYZ", 1e3),
	strings.Repeat("\xff\x00", 1e4) + "hotdog",
}

func TestRolling(t *testing.T) {
	blocksize := 16
	rolling := New()

	for _, teststring := range rollStrings {
		test := []byte(teststring)
		for i := 0; i < len(test)-blocksize-1; i++ {
			in := test[i : i+blocksize]
			if i == 0 {
				rolling.Write(in)
			}

			slow := checksum(in)
			if got := rolling.Sum32(); got != slow {
				t.Errorf("rolling hash: %q = %d want %d", in, got, slow)
			}

			rolling.Roll(uint32(blocksize), test[i], test[i+blocksize])
		}
		rolling.Reset()
	}
}
