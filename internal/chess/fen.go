package chess

type FEN struct {
	BoardState      string
	ActiveColour    string
	CastlingRights  string
	EnPassantSquare string
	HalfMoveClock   string
	FullMoveNumber  string
}

func (f FEN) toString() string {
	return f.BoardState + " " + f.ActiveColour + " " + f.CastlingRights + " " + f.EnPassantSquare + " " + f.HalfMoveClock + " " + f.FullMoveNumber
}
