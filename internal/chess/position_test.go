package chess

import "testing"

func TestLoadFEN(t *testing.T) {
	type squareCheck struct {
		coordinate
		squareState
	}

	type testSpec struct {
		name string
		FEN
		squareChecks    []squareCheck
		enPassantSquare coordinate
		castlingRights
		activeColour  colour
		halfMoveClock uint8
	}

	tests := []testSpec{
		{
			"opera",
			FEN{
				BoardState:      "3rkb1r/p2nqppp/5n2/1B2p1B1/4P3/1Q6/PPP2PPP/2KR3R",
				ActiveColour:    "w",
				CastlingRights:  "k",
				EnPassantSquare: "-",
				HalfMoveClock:   "3",
				FullMoveNumber:  "13",
			},
			[]squareCheck{
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

	runTest := func(t *testing.T, test testSpec) {
		thePosition := Position{}
		thePosition.LoadFEN(test.FEN)

		for _, sc := range test.squareChecks {
			if result := thePosition.board[sc.coordinate]; result != sc.squareState {
				t.Errorf("expected %v at %v, got %v", sc.squareState, sc.coordinate, result)
			}

			if sc.squareState != empty {
				theColour := sc.squareState.getColour()
				thePieceType := sc.squareState.getPieceType()
				if !thePosition.colourMasks[theColour].get(sc.coordinate) {
					t.Errorf("colourMask not set for %v at %v", sc.squareState, sc.coordinate)
				}
				if !thePosition.pieceTypeMasks[thePieceType].get(sc.coordinate) {
					t.Errorf("pieceTypeMask not set for %v at %v", sc.squareState, sc.coordinate)
				}
			}
		}

		if thePosition.enPassantSquare != test.enPassantSquare {
			t.Errorf("expected en passant square %v, got %v", test.enPassantSquare, thePosition.enPassantSquare)
		}

		if thePosition.castlingRights != test.castlingRights {
			t.Errorf("expected castling rights %v, got %v", test.castlingRights, thePosition.castlingRights)
		}

		if thePosition.activeColour != test.activeColour {
			t.Errorf("should be %v to move", test.activeColour)
		}

		if thePosition.halfMoveClock != test.halfMoveClock {
			t.Errorf("expected half move clock to be %v, got %v", test.halfMoveClock, thePosition.halfMoveClock)
		}
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) { runTest(t, test) })
	}
}
