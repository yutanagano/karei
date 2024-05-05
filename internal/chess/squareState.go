package chess

import "fmt"

type squareState uint8

const (
	whiteKing squareState = iota
	blackKing
	whiteQueen
	blackQueen
	whiteRook
	blackRook
	whiteBishop
	blackBishop
	whiteKnight
	blackKnight
	whitePawn
	blackPawn
	empty
)

var squareStateRunes = "KkQqRrBbNnPp "

var runeToSquareStateMap = map[rune]squareState{
	'K': whiteKing,
	'Q': whiteQueen,
	'R': whiteRook,
	'B': whiteBishop,
	'N': whiteKnight,
	'P': whitePawn,
	'k': blackKing,
	'q': blackQueen,
	'r': blackRook,
	'b': blackBishop,
	'n': blackKnight,
	'p': blackPawn,
}

func squareStateFromRune(char rune) (squareState, error) {
	s, ok := runeToSquareStateMap[char]

	if !ok {
		return s, fmt.Errorf("unrecognised square state %c", char)
	}

	return s, nil
}

func (s squareState) toRune() rune {
	return rune(squareStateRunes[s])
}

func (s squareState) getPieceType() pieceType {
	switch s / 6 {
	case 0:
		return king
	case 1:
		return queen
	case 2:
		return rook
	case 3:
		return bishop
	case 4:
		return knight
	default:
		return pawn
	}
}

func (s squareState) getColour() colour {
	if s%2 == 0 {
		return white
	}
	return black
}
