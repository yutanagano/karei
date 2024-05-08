package chess

type colour uint8
type pieceType uint8

const (
	white colour = 0
	black colour = 1
)

const (
	king pieceType = iota
	queen
	rook
	bishop
	knight
	pawn
)

func (c colour) getOpponent() colour {
	if c == white {
		return black
	}
	return white
}
