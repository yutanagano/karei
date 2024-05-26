package chess

import (
	"fmt"
	"unicode"
)

type pieceType uint8

const (
	king pieceType = iota
	queen
	rook
	bishop
	knight
	pawn
)

type colour uint8

const (
	white colour = iota
	black
)

func (c colour) getOpponent() colour {
	if c == white {
		return black
	}
	return white
}

type piece struct {
	colour
	pieceType
}

func pieceFromRune(r rune) (piece, error) {
	var c colour
	var p pieceType

	if unicode.IsUpper(r) {
		c = white
	} else {
		c = black
	}

	switch unicode.ToLower(r) {
	case 'k':
		p = king
	case 'q':
		p = queen
	case 'r':
		p = rook
	case 'b':
		p = bishop
	case 'p':
		p = pawn
	default:
		return piece{}, fmt.Errorf("rune does not correspond to piece: %c", r)
	}

	return piece{c, p}, nil
}

func (p piece) getColour() colour {
	return p.colour
}

func (p piece) getPieceType() pieceType {
	return p.pieceType
}
