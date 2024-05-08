package chess

import (
	"testing"
)

var operaGame = FEN{
	BoardState:      "3rkb1r/p2nqppp/5n2/1B2p1B1/4P3/1Q6/PPP2PPP/2KR3R",
	ActiveColour:    "w",
	CastlingRights:  "k",
	EnPassantSquare: "-",
	HalfMoveClock:   "3",
	FullMoveNumber:  "13",
}

func TestLoadFEN(t *testing.T) {
	type testCase struct {
		name string
		FEN
		squareChecks []struct {
			coordinate
			squareState
		}
		enPassantSquare coordinate
		castlingRights
		activeColour  colour
		halfMoveClock uint8
	}

	testCases := []testCase{
		{
			"opera",
			operaGame,
			[]struct {
				coordinate
				squareState
			}{
				{d8, blackRook},
				{c1, whiteKing},
				{g7, blackPawn},
				{g5, whiteBishop},
				{e1, empty},
			},
			nullCoordinate,
			0b0100,
			white,
			3,
		},
	}

	checkCase := func(t *testing.T, c testCase) {
		thePosition := Position{}
		thePosition.LoadFEN(c.FEN)

		for _, sc := range c.squareChecks {
			if result := thePosition.board[sc.coordinate]; result != sc.squareState {
				t.Errorf("expected %v at %v, got %v", sc.squareState, sc.coordinate, result)
			}

			if sc.squareState != empty {
				theColour := sc.squareState.getColour()
				thePieceType := sc.squareState.getPieceType()
				if !thePosition.occupationByColour[theColour].get(sc.coordinate) {
					t.Errorf("colourMask not set for %v at %v", sc.squareState, sc.coordinate)
				}
				if !thePosition.occupationByPieceType[thePieceType].get(sc.coordinate) {
					t.Errorf("pieceTypeMask not set for %v at %v", sc.squareState, sc.coordinate)
				}
			}
		}

		if thePosition.enPassantSquare != c.enPassantSquare {
			t.Errorf("expected en passant square %v, got %v", c.enPassantSquare, thePosition.enPassantSquare)
		}

		if thePosition.castlingRights != c.castlingRights {
			t.Errorf("expected castling rights %v, got %v", c.castlingRights, thePosition.castlingRights)
		}

		if thePosition.activeColour != c.activeColour {
			t.Errorf("should be %v to move", c.activeColour)
		}

		if thePosition.halfMoveClock != c.halfMoveClock {
			t.Errorf("expected half move clock to be %v, got %v", c.halfMoveClock, thePosition.halfMoveClock)
		}
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) { checkCase(t, c) })
	}
}

func TestGetPseudoLegalMoves(t *testing.T) {
	// TODO test from another position to look at knights, pawns, castling and en-passant
	thePosition := Position{}
	thePosition.LoadFEN(operaGame)

	expectedNumMoves := 41
	expected := moveList{
		{d1, d6, empty},
		{d1, f1, empty},
		{h1, g1, empty},
		{b5, d7, empty},
		{g5, h4, empty},
		{b3, h3, empty},
		{b3, f7, empty},
		{c1, d2, empty},
	}

	if numMoves := len(thePosition.legalMoves); numMoves != expectedNumMoves {
		t.Errorf("expected %v moves, got %v", expectedNumMoves, numMoves)
	}

	for _, m := range expected {
		if !thePosition.legalMoves.contains(m) {
			t.Errorf("expected move %v that was not generated", m)
		}
	}
}
