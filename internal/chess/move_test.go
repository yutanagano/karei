package chess

import (
	"reflect"
	"testing"
)

func TestMoveFromString(t *testing.T) {
	type testCase struct {
		moveString string
		expected   move
	}

	testCases := []testCase{
		{
			"e2e4",
			move{e2, e4, empty},
		},
		{
			"f7f8Q",
			move{f7, f8, whiteQueen},
		},
	}

	checkCase := func(t *testing.T, c testCase) {
		result, err := moveFromString(c.moveString)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(result, c.expected) {
			t.Errorf("expected %v, got %v", c.expected, result)
		}
	}

	for _, c := range testCases {
		t.Run(c.moveString, func(t *testing.T) { checkCase(t, c) })
	}
}

func TestMoveToString(t *testing.T) {
	type testCase struct {
		move     move
		expected string
	}

	testCases := []testCase{
		{
			move{d7, d5, empty},
			"d7d5",
		},
		{
			move{h2, h1, blackQueen},
			"h2h1q",
		},
	}

	checkCase := func(t *testing.T, c testCase) {
		result := c.move.toString()
		if result != c.expected {
			t.Errorf("expected %s, got %s", c.expected, result)
		}
	}

	for _, c := range testCases {
		t.Run(c.expected, func(t *testing.T) { checkCase(t, c) })
	}
}

func TestMoveListFilter(t *testing.T) {
	type testCase struct {
		name         string
		initialList  moveList
		evaluator    func(move) bool
		expectedList moveList
	}

	testCases := []testCase{
		{
			"non-promotions",
			moveList{
				{e2, e4, empty},
				{d2, d4, empty},
				{e7, e8, whiteQueen},
				{d7, d8, whiteQueen},
			},
			func(theMove move) bool {
				return theMove.Promotion == empty
			},
			moveList{
				{e2, e4, empty},
				{d2, d4, empty},
			},
		},
		{
			"from e2",
			moveList{
				{e2, e3, empty},
				{e2, e4, empty},
				{d2, d3, empty},
				{d2, d4, empty},
			},
			func(theMove move) bool {
				return theMove.From == e2
			},
			moveList{
				{e2, e3, empty},
				{e2, e4, empty},
			},
		},
	}

	checkCase := func(t *testing.T, c testCase) {
		c.initialList.filter(c.evaluator)
		if !reflect.DeepEqual(c.initialList, c.expectedList) {
			t.Errorf("expected %v, got %v", c.expectedList, c.initialList)
		}
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) { checkCase(t, c) })
	}
}
