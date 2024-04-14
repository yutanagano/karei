package main

import "strconv"

type Fen struct {
	BoardState      string
	ActiveColor     string
	CastlingRights  string
	EnPassantSquare string
	HalfMoveClock   int
	FullMoveNumber  int
}

func (f Fen) ToString() string {
	return f.BoardState + " " + f.ActiveColor + " " + f.CastlingRights + " " + f.EnPassantSquare + " " + strconv.Itoa(f.HalfMoveClock) + " " + strconv.Itoa(f.FullMoveNumber)
}

func parseFen(f Fen) {
}

func parseMoves(tokens *[]string) {
}
