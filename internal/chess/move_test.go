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

func TestToString(t *testing.T) {
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
