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

func GetStartingFEN() FEN {
	return FEN{
		BoardState:      "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR",
		ActiveColour:    "w",
		CastlingRights:  "KQkq",
		EnPassantSquare: "-",
		HalfMoveClock:   "0",
		FullMoveNumber:  "1",
	}
}
