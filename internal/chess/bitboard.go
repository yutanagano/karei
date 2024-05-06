package chess

import (
	"math/bits"
)

type bitBoard uint64

func (b bitBoard) count() int {
	return bits.OnesCount64(uint64(b))
}

func (b bitBoard) get(coord coordinate) bool {
	return b&(1<<coord) != 0
}

func (b *bitBoard) turnOn(coord coordinate) {
	*b |= 1 << coord
}

func (b *bitBoard) turnOff(coord coordinate) {
	*b &= bitBoard(^uint64(1) << coord)
}

func (b *bitBoard) pop() (place coordinate, ok bool) {
	if *b == 0 {
		place = 0
		ok = false
		return place, ok
	}

	place = coordinate(bits.TrailingZeros64(uint64(*b)))
	ok = true
	*b &= ^(1 << place)

	return place, ok
}
