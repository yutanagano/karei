package chess

import (
	"reflect"
	"testing"
)

func TestAlgebraicMoveFromString(t *testing.T) {
	type testCase struct {
		moveString string
		expected   algebraicMove
	}

	testCases := []testCase{
		{
			"e2e4",
			algebraicMove{e2, e4, empty},
		},
		{
			"f7f8Q",
			algebraicMove{f7, f8, whiteQueen},
		},
	}

	checkCase := func(t *testing.T, c testCase) {
		result, err := algebraicMoveFromString(c.moveString)
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

func TestAlgebraicMoveToString(t *testing.T) {
	type testCase struct {
		move     algebraicMove
		expected string
	}

	testCases := []testCase{
		{
			algebraicMove{d7, d5, empty},
			"d7d5",
		},
		{
			algebraicMove{h2, h1, blackQueen},
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

func TestMoveFromParts(t *testing.T) {
	type testCase struct {
		name                  string
		from                  coordinate
		to                    coordinate
		capturedPiece         squareState
		promotionTo           squareState
		currentCastlingRights castlingRights
		currentEPSquare       coordinate
		expectedMove          move
	}

	testCases := []testCase{
		{
			"e2e4",
			e2, e4, empty, empty, castlingRights(0b1111), nullCoordinate,
			move(0x40fcc70c),
		},
	}

	checkCase := func(t *testing.T, c testCase) {
		result := moveFromParts(c.from, c.to, c.capturedPiece, c.promotionTo, c.currentCastlingRights, c.currentEPSquare)
		if result != c.expectedMove {
			t.Errorf("expected %v, got %v", c.expectedMove, result)
		}
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) { checkCase(t, c) })
	}
}

func TestMoveGetParts(t *testing.T) {
	theMove := moveFromParts(h7, g8, blackQueen, whiteQueen, blackCastleKingSide|blackCastleQueenSide, nullCoordinate)

	if result := theMove.getFromCoordinate(); result != h7 {
		t.Errorf("expected h7, got %v", result)
	}

	if result := theMove.getToCoordinate(); result != g8 {
		t.Errorf("expected g8, got %v", result)
	}

	if result := theMove.getCapturedPiece(); result != blackQueen {
		t.Errorf("expected blackQueen, got %v", result)
	}

	if result := theMove.getPromotionTo(); result != whiteQueen {
		t.Errorf("expected whiteQueen, got %v", result)
	}

	if result := theMove.getCurrentCastlingRights(); result != 0b1100 {
		t.Errorf("expected castling rights of kq, got %v", result)
	}

	if result := theMove.getCurrentEPSquare(); result != nullCoordinate {
		t.Errorf("expected nullCoordinate, got %v", result)
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
			moveList{0, 1, 2, 3, 4, 5},
			func(theMove move) bool {
				return theMove%2 == 0
			},
			moveList{0, 2, 4},
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
