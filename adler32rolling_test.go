package adler32rolling

import (
	//	"fmt"
	"strings"
	"testing"
)

var golden = []struct {
	out uint32
	in  string
}{
	{0x00000001, ""},
	{0x5f1007e3, "You're a lizard, Harry."},
}

// checksum is a slow but simple implementation of the Adler-32 checksum.
// It is a straight port of the sample code in RFC 1950 section 9.
func checksum(p []byte) uint32 {
	s1, s2 := uint32(1), uint32(0)
	for _, x := range p {
		s1 = (s1 + uint32(x)) % mod
		s2 = (s2 + s1) % mod
	}
	return s2<<16 | s1
}

func TestGolden(t *testing.T) {
	for _, g := range golden {
		in := g.in
		p := []byte(g.in)
		if got := checksum(p); got != g.out {
			t.Errorf("simple implementation: checksum(%q) = 0x%x want 0x%x", in, got, g.out)
			continue
		}

		if got := Checksum(p); got != g.out {
			t.Errorf("optimized implementation: Checksum(%q) = 0x%x want 0x%x", in, got, g.out)
			continue
		}
	}
}

var rollStrings = []string{
	"You're a lizard, Harry.",
	strings.Repeat("ABCDEFGHIJKLMNOPQRSTUVWXYZ", 1e4),
}

func TestRolling(t *testing.T) {
	blocksize := 16
	//	rollingString := []byte("You're a lizard, Harry.")

	hash := New()
	rolling := New()

	for _, teststring := range rollStrings {
		test := []byte(teststring)
		for i := 0; i < len(test)-blocksize-1; i++ {
			in := test[i : i+blocksize]
			hash.Write(in)
			if i == 0 {
				rolling.Write(in)
			}
			//fmt.Printf("%s\t %d %d\n", in, hash.Sum32(), i)

			if got := rolling.Sum32(); got != hash.Sum32() {
				t.Errorf("rolling hash: %q = %d want %d", in, got, hash.Sum32())
			}

			rolling.Roll(uint32(blocksize), test[i+blocksize], test[i])

			hash.Reset()
		}
		rolling.Reset()
	}

}
