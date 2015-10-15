// This is practically entirely a copy of the hash/adler32/adler32.go
// except it also includes a Roll() function that allows to pop an old
// byte and push a new byte, giving a rolling hash method.

package adler32rolling

import "hash"

const (
	mod  = 65521
	nmax = 5552
)

const Size = 4

type digest uint32

func (d *digest) Reset() { *d = 0 }

type Hash32 interface {
	hash.Hash32
	Roll(blocksize uint32, p byte, r byte)
}

func New() Hash32 {
	d := new(digest)
	d.Reset()
	return d
}

func (d *digest) Size() int { return Size }

func (d *digest) BlockSize() int { return 1 }

// Add p to the running checksum d.
func update(d digest, p []byte) digest {
	s1, s2 := uint32(d&0xffff), uint32(d>>16)
	for len(p) > 0 {
		var q []byte
		if len(p) > nmax {
			p, q = p[:nmax], p[nmax:]
		}
		for _, x := range p {
			s1 += uint32(x)
			s2 += s1
		}
		s1 %= mod
		s2 %= mod
		p = q
	}
	return digest(s2<<16 | s1)
}

func (d *digest) Write(p []byte) (nn int, err error) {
	*d = update(*d, p)
	return len(p), nil
}

// Add p to the running checksum d while removing r.
func roll(d digest, blocksize uint32, in byte, out byte) digest {
	s1, s2 := uint32(d&0xffff), uint32(d>>16)

	i := uint32(in)
	o := uint32(out)

	/*
		s1 -= o
		s2 -= blocksize * o

		s1 += i
		s2 += s1
	*/

	s1 += i - o
	s2 += s1 - blocksize*o

	s1 %= mod
	s2 %= mod

	return digest(s2<<16 | s1)
}

func (d *digest) Roll(blocksize uint32, p byte, r byte) {
	*d = roll(*d, blocksize, p, r)
}

func (d *digest) Sum32() uint32 { return uint32(*d) }

func (d *digest) Sum(in []byte) []byte {
	s := uint32(*d)
	return append(in, byte(s>>24), byte(s>>16), byte(s>>8), byte(s))
}

func Checksum(data []byte) uint32 { return uint32(update(1, data)) }
