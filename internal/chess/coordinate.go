package chess

import "fmt"

type coordinate uint8

const (
	a1 coordinate = iota
	b1
	c1
	d1
	e1
	f1
	g1
	h1
	a2
	b2
	c2
	d2
	e2
	f2
	g2
	h2
	a3
	b3
	c3
	d3
	e3
	f3
	g3
	h3
	a4
	b4
	c4
	d4
	e4
	f4
	g4
	h4
	a5
	b5
	c5
	d5
	e5
	f5
	g5
	h5
	a6
	b6
	c6
	d6
	e6
	f6
	g6
	h6
	a7
	b7
	c7
	d7
	e7
	f7
	g7
	h7
	a8
	b8
	c8
	d8
	e8
	f8
	g8
	h8
	nullCoordinate
)

var stringToCoordMap = map[string]coordinate{
	"a1": a1,
	"b1": b1,
	"c1": c1,
	"d1": d1,
	"e1": e1,
	"f1": f1,
	"g1": g1,
	"h1": h1,
	"a2": a2,
	"b2": b2,
	"c2": c2,
	"d2": d2,
	"e2": e2,
	"f2": f2,
	"g2": g2,
	"h2": h2,
	"a3": a3,
	"b3": b3,
	"c3": c3,
	"d3": d3,
	"e3": e3,
	"f3": f3,
	"g3": g3,
	"h3": h3,
	"a4": a4,
	"b4": b4,
	"c4": c4,
	"d4": d4,
	"e4": e4,
	"f4": f4,
	"g4": g4,
	"h4": h4,
	"a5": a5,
	"b5": b5,
	"c5": c5,
	"d5": d5,
	"e5": e5,
	"f5": f5,
	"g5": g5,
	"h5": h5,
	"a6": a6,
	"b6": b6,
	"c6": c6,
	"d6": d6,
	"e6": e6,
	"f6": f6,
	"g6": g6,
	"h6": h6,
	"a7": a7,
	"b7": b7,
	"c7": c7,
	"d7": d7,
	"e7": e7,
	"f7": f7,
	"g7": g7,
	"h7": h7,
	"a8": a8,
	"b8": b8,
	"c8": c8,
	"d8": d8,
	"e8": e8,
	"f8": f8,
	"g8": g8,
	"h8": h8,
}

var coordToStringMap = map[coordinate]string{
	a1: "a1",
	b1: "b1",
	c1: "c1",
	d1: "d1",
	e1: "e1",
	f1: "f1",
	g1: "g1",
	h1: "h1",
	a2: "a2",
	b2: "b2",
	c2: "c2",
	d2: "d2",
	e2: "e2",
	f2: "f2",
	g2: "g2",
	h2: "h2",
	a3: "a3",
	b3: "b3",
	c3: "c3",
	d3: "d3",
	e3: "e3",
	f3: "f3",
	g3: "g3",
	h3: "h3",
	a4: "a4",
	b4: "b4",
	c4: "c4",
	d4: "d4",
	e4: "e4",
	f4: "f4",
	g4: "g4",
	h4: "h4",
	a5: "a5",
	b5: "b5",
	c5: "c5",
	d5: "d5",
	e5: "e5",
	f5: "f5",
	g5: "g5",
	h5: "h5",
	a6: "a6",
	b6: "b6",
	c6: "c6",
	d6: "d6",
	e6: "e6",
	f6: "f6",
	g6: "g6",
	h6: "h6",
	a7: "a7",
	b7: "b7",
	c7: "c7",
	d7: "d7",
	e7: "e7",
	f7: "f7",
	g7: "g7",
	h7: "h7",
	a8: "a8",
	b8: "b8",
	c8: "c8",
	d8: "d8",
	e8: "e8",
	f8: "f8",
	g8: "g8",
	h8: "h8",
}

func coordinateFromRankFileIndices(rankIndex int8, fileIndex int8) (coordinate, error) {
	coordinateInt := rankIndex*8 + fileIndex

	if coordinateInt < 0 || coordinateInt >= 64 {
		return nullCoordinate, fmt.Errorf("invalid rank/file: %v/%v", rankIndex, fileIndex)
	}

	return coordinate(coordinateInt), nil
}

func coordinateFromString(s string) (coordinate, error) {
	c, ok := stringToCoordMap[s]

	if ok {
		return c, nil
	}

	return c, fmt.Errorf("Unrecognised coordinate name %s", s)
}

func (c coordinate) toString() string {
	return coordToStringMap[c]
}

func (c coordinate) getRankIndex() int8 {
	return int8(c) / 8
}

func (c coordinate) getFileIndex() int8 {
	return int8(c) % 8
}
