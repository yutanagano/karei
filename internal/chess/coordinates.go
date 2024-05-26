package chess

import (
	"fmt"
	"strconv"
)

const (
	fileA uint8 = iota
	fileB
	fileC
	fileD
	fileE
	fileF
	fileG
	fileH
)

const (
	rank1 uint8 = iota
	rank2
	rank3
	rank4
	rank5
	rank6
	rank7
	rank8
)

type coordinate uint8

const (
	a1 coordinate = iota
	a2
	a3
	a4
	a5
	a6
	a7
	a8
	b1
	b2
	b3
	b4
	b5
	b6
	b7
	b8
	c1
	c2
	c3
	c4
	c5
	c6
	c7
	c8
	d1
	d2
	d3
	d4
	d5
	d6
	d7
	d8
	e1
	e2
	e3
	e4
	e5
	e6
	e7
	e8
	f1
	f2
	f3
	f4
	f5
	f6
	f7
	f8
	g1
	g2
	g3
	g4
	g5
	g6
	g7
	g8
	h1
	h2
	h3
	h4
	h5
	h6
	h7
	h8
)

type gridOffset struct{ fileOffset, rankOffset int8 }

func coordinateFromParts(file, rank uint8) coordinate {
	return coordinate(file + rank*8)
}

func coordinateFromString(s string) (coordinate, error) {
	var fileIdx, rankIdx uint8
	potentialError := fmt.Errorf("cannot convert string to coordinate: %s", s)

	if len(s) != 2 {
		return coordinate(0), potentialError
	}

	switch s[0] {
	case 'a':
		fileIdx = 0
	case 'b':
		fileIdx = 1
	case 'c':
		fileIdx = 2
	case 'd':
		fileIdx = 3
	case 'e':
		fileIdx = 4
	case 'f':
		fileIdx = 5
	case 'g':
		fileIdx = 6
	case 'h':
		fileIdx = 7
	default:
		return coordinate(0), potentialError
	}

	rankNumber, err := strconv.Atoi(string(s[1]))
	if err != nil {
		return coordinate(0), potentialError
	}
	if rankNumber == 0 {
		return coordinate(0), potentialError
	} else if rankNumber > 8 {
		return coordinate(0), potentialError
	}
	rankIdx = uint8(rankNumber - 1)

	return coordinate(fileIdx + rankIdx*8), nil
}

func (c coordinate) getFile() uint8 {
	return uint8(c % 8)
}

func (c coordinate) getRank() uint8 {
	return uint8(c / 8)
}

func (c *coordinate) move(o gridOffset) (coordinate, error) {
	file := int8(*c%8) + o.fileOffset
	if file < 0 {
		return coordinate(0), fmt.Errorf("file out of bounds.")
	} else if file > 7 {
		return coordinate(0), fmt.Errorf("file out of bounds.")
	}

	rank := int8(*c/8) + o.rankOffset
	if rank < 0 {
		return coordinate(0), fmt.Errorf("rank out of bounds.")
	} else if rank > 7 {
		return coordinate(0), fmt.Errorf("rank out of bounds.")
	}

	return coordinateFromParts(uint8(file), uint8(rank)), nil
}
