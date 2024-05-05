package chess

import (
	"reflect"
	"testing"
)

func TestMoveFromString(t *testing.T) {
	type testCase struct {
		moveString string
		expected   Move
	}

	testCases := []testCase{
		{
			"e2e4",
			Move{e2, e4, empty},
		},
		{
			"f7f8Q",
			Move{f7, f8, whiteQueen},
		},
	}

	checkCase := func(t *testing.T, c testCase) {
		result, err := MoveFromString(c.moveString)
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
		move     Move
		expected string
	}

	testCases := []testCase{
		{
			Move{d7, d5, empty},
			"d7d5",
		},
		{
			Move{h2, h1, blackQueen},
			"h2h1q",
		},
	}

	checkCase := func(t *testing.T, c testCase) {
		result := c.move.ToString()
		if result != c.expected {
			t.Errorf("expected %s, got %s", c.expected, result)
		}
	}

	for _, c := range testCases {
		t.Run(c.expected, func(t *testing.T) { checkCase(t, c) })
	}
}
