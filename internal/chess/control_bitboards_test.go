package chess

import "testing"

func TestKingControlBitBoards(t *testing.T) {
	type testCase struct {
		currentSquare coordinate
		expected      bitBoard
	}

	testCases := []testCase{
		{a1, bitBoard(0b1100000010)},
		{e4, bitBoard(0b11100000101000001110000000000000000000)},
	}

	for _, c := range testCases {
		result := kingControlFrom[c.currentSquare]
		if result != c.expected {
			t.Errorf("expected %v, got %v", c.expected, result)
		}
	}
}

func TestKnightControlBitBoards(t *testing.T) {
	type testCase struct {
		currentSquare coordinate
		expected      bitBoard
	}

	testCases := []testCase{
		{f3, bitBoard(0b101000010001000000000001000100001010000)},
		{b5, bitBoard(0b101000010000000000000001000000001010000000000000000)},
	}

	for _, c := range testCases {
		result := knightControlFrom[c.currentSquare]
		if result != c.expected {
			t.Errorf("expected %v, got %v", c.expected, result)
		}
	}
}
