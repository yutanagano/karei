package chess

import (
	"math/bits"
)

type bitBoard uint64

func (b bitBoard) count() int {
	return bits.OnesCount64(uint64(b))
}

func (b *bitBoard) set(coord coordinate) {
	*b |= 1 << coord
}

func (b bitBoard) get(coord coordinate) bool {
	return b&(1<<coord) != 0
}

func (b *bitBoard) clear(coord coordinate) {
	*b &= bitBoard(^uint64(1) << coord)
}
