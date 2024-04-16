package main

type bitBoard uint64
type colour uint8
type castlingRights uint8
type pieceType uint8
type squareState uint8

const (
	a1 = iota
	b1
	c1
	d1
	e1
	f1
	g1
	h1
	a2
	b2
	c2
	d2
	e2
	f2
	g2
	h2
	a3
	b3
	c3
	d3
	e3
	f3
	g3
	h3
	a4
	b4
	c4
	d4
	e4
	f4
	g4
	h4
	a5
	b5
	c5
	d5
	e5
	f5
	g5
	h5
	a6
	b6
	c6
	d6
	e6
	f6
	g6
	h6
	a7
	b7
	c7
	d7
	e7
	f7
	g7
	h7
	a8
	b8
	c8
	d8
	e8
	f8
	g8
	h8
)

const (
	white colour = 0
	black colour = 1
)

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

const (
	king pieceType = iota
	queen
	rook
	bishop
	knight
	pawn
)

const (
	whiteCastleKingSide  castlingRights = 0b0001
	whiteCastleQueenSide castlingRights = 0b0010
	blackCastleKingSide  castlingRights = 0b0100
	blackCastleQueenSide castlingRights = 0b1000
)

type fen struct {
	boardState      string
	activeColour    string
	castlingRights  string
	enPassantSquare string
	halfMoveClock   string
	fullMoveNumber  string
}

func (f fen) toString() string {
	return f.boardState + " " + f.activeColour + " " + f.castlingRights + " " + f.enPassantSquare + " " + f.halfMoveClock + " " + f.fullMoveNumber
}

func (f fen) toPosition() {
	// TODO: Return position struct from FEN
}

type position struct {
	board           [64]squareState
	colorMasks      [2]bitBoard
	typeMasks       [6]bitBoard
	kingSquares     [2]int
	enPassantSquare int
	castlingRights  castlingRights
	activeColour    colour
	pieceCounter    [12]int
	halfMoveClock   int
}

func (p position) pieceMask() bitBoard {
	return p.colorMasks[white] | p.colorMasks[black]
}

func (p *position) clear() {
	for coordinate := a1; coordinate <= h8; coordinate++ {
		p.board[coordinate] = empty
	}

	p.colorMasks[white] = 0
	p.colorMasks[black] = 0

	for piece := king; piece <= pawn; piece++ {
		p.typeMasks[piece] = 0
	}

	p.kingSquares[white] = e1
	p.kingSquares[black] = e8

	p.enPassantSquare = -1
	p.castlingRights = whiteCastleKingSide | whiteCastleQueenSide | blackCastleKingSide | blackCastleQueenSide
	p.activeColour = white
	for piece := whiteKing; piece <= blackPawn; piece++ {
		p.pieceCounter[piece] = 0
	}
}

func (p *position) newGame() {
	p.clear()
	// TODO: set up starting position
}

func parseMoves(tokens *[]string) {
}
