package chess

import (
	"strconv"
	"testing"
)

func TestSet(t *testing.T) {
	theBitBoard := bitBoard(0)
	theBitBoard.set(e4)
	expected := bitBoard(1 << e4)

	if theBitBoard != expected {
		t.Errorf("expected %v, got %v", expected, theBitBoard)
	}
}

func TestGet(t *testing.T) {
	theBitBoard := bitBoard(1 << c6)

	if theBitBoard.get(c6) != true {
		t.Errorf("the c6 square should be set")
	}

	if theBitBoard.get(h3) != false {
		t.Errorf("the h3 square should not be set")
	}
}

func TestClear(t *testing.T) {
	theBitBoard := bitBoard(1 << d5)
	theBitBoard.clear(d5)

	if theBitBoard != 0 {
		t.Errorf("the d5 square should be cleared")
	}
}

func TestPop(t *testing.T) {
	type testCase struct {
		originalBitBoard bitBoard
		expectedOk       bool
		expectedPlace    int
		expectedBitBoard bitBoard
	}

	testCases := []testCase{
		{bitBoard(0b010101), true, 0, bitBoard(0b010100)},
		{bitBoard(0b011000), true, 3, bitBoard(0b010000)},
		{bitBoard(0), false, 0, 0},
	}

	checkCase := func(t *testing.T, c testCase) {
		place, ok := c.originalBitBoard.pop()
		if ok != c.expectedOk {
			t.Errorf("expected ok of %v", c.expectedOk)
		}
		if c.expectedOk == false {
			return
		}

		if place != c.expectedPlace {
			t.Errorf("expected place %v, got %v", c.expectedPlace, place)
		}
		if c.originalBitBoard != c.expectedBitBoard {
			t.Errorf("expected %v after popping, got %v", c.expectedBitBoard, c.originalBitBoard)
		}
	}

	for idx, c := range testCases {
		t.Run(strconv.Itoa(idx), func(t *testing.T) { checkCase(t, c) })
	}
}
