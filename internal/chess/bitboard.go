package chess

import (
	"math/bits"
)

type bitBoard uint64

const (
	bitBoardFileAMask bitBoard = 0x0101010101010101
	bitBoardFileHMask bitBoard = 0x8080808080808080
	bitBoardRank4Mask bitBoard = 0x00000000FF000000
	bitBoardRank5Mask bitBoard = 0x000000FF00000000
)

func (b bitBoard) count() int {
	return bits.OnesCount64(uint64(b))
}

func (b bitBoard) get(c coordinate) bool {
	return b&(1<<c) != 0
}

func (b *bitBoard) turnOn(c coordinate) {
	*b |= 1 << c
}

func (b *bitBoard) turnOff(c coordinate) {
	*b &= ^(1 << c)
}

func (b *bitBoard) pop() (c coordinate, ok bool) {
	if *b == 0 {
		return coordinate(0), false
	}

	c = coordinate(bits.TrailingZeros64(uint64(*b)))
	*b &= ^(1 << c)
	return c, true
}
