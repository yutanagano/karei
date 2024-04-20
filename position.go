package main

import (
	"fmt"
	"strconv"
)

type colour uint8
type castlingRights uint8
type pieceType uint8
type squareState uint8

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

type position struct {
	board            [64]squareState
	colourMasks      [2]bitBoard
	pieceTypeMasks   [6]bitBoard
	kingSquares      [2]coordinate
	enPassantSquare  int
	castlingRights   castlingRights
	activeColour     colour
	pieceTypeCounter [12]int
	halfMoveClock    int
}

func (p position) occupiedMask() bitBoard {
	return p.colourMasks[white] | p.colourMasks[black]
}

func (p *position) clear() {
	for theCoord := a1; theCoord <= h8; theCoord++ {
		p.board[theCoord] = empty
	}

	p.colourMasks[white] = 0
	p.colourMasks[black] = 0

	for thePieceType := king; thePieceType <= pawn; thePieceType++ {
		p.pieceTypeMasks[thePieceType] = 0
	}

	p.kingSquares[white] = e1
	p.kingSquares[black] = e8

	p.enPassantSquare = -1
	p.castlingRights = whiteCastleKingSide | whiteCastleQueenSide | blackCastleKingSide | blackCastleQueenSide
	p.activeColour = white
	for thePiece := whiteKing; thePiece <= blackPawn; thePiece++ {
		p.pieceTypeCounter[thePiece] = 0
	}
}

func (p *position) loadFEN(f fen) error {
	p.clear()

	currentRow, currentColumn := 7, 0
	for _, currentRune := range f.boardState {
		if currentColumn > 8 {
			err := fmt.Errorf("Bad FEN: overfilled row during board specification.")
			return err
		}

		if currentRune == '/' {
			if currentColumn != 8 {
				err := fmt.Errorf("Bad FEN: underfilled row during board specification.")
				return err
			}

			currentRow--
			currentColumn = 0
			continue
		}

		if numEmptySquares, err := strconv.Atoi(string(currentRune)); err != nil {
			currentColumn += numEmptySquares
			continue
		}

		currentCoordinate := coordinateFromRowColumn(currentRow, currentColumn)
		currentSquareState, err := squareStateFromRune(currentRune)

		if err != nil {
			err = fmt.Errorf("Bad FEN: %s", err.Error())
			return err
		}

		p.setSquare(currentSquareState, currentCoordinate)
	}

	return nil
}

func (p *position) setSquare(state squareState, coord coordinate) {
	p.board[coord] = state

	if state == empty {
		p.colourMasks[white].clear(coord)
		p.colourMasks[black].clear(coord)

		for thePieceType := king; thePieceType <= pawn; thePieceType++ {
			p.pieceTypeMasks[thePieceType].clear(coord)
		}
		return
	}

	thePieceType := state.getPieceType()
	theColour := state.getColour()

	if thePieceType == king {
		p.kingSquares[theColour] = coord
	}

	p.colourMasks[theColour].set(coord)
	p.pieceTypeMasks[thePieceType].set(coord)
}

func (p *position) newGame() {
	p.clear()
	// TODO: set up starting position
}

func squareStateFromRune(char rune) (squareState, error) {
	switch char {
	case 'K':
		return whiteKing, nil
	case 'Q':
		return whiteQueen, nil
	case 'R':
		return whiteRook, nil
	case 'B':
		return whiteBishop, nil
	case 'N':
		return whiteKnight, nil
	case 'P':
		return whitePawn, nil
	case 'k':
		return blackKing, nil
	case 'q':
		return blackQueen, nil
	case 'r':
		return blackRook, nil
	case 'b':
		return blackBishop, nil
	case 'n':
		return blackKnight, nil
	case 'p':
		return blackPawn, nil
	default:
		err := fmt.Errorf("Unrecognized piece: %c", char)
		return empty, err
	}
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

func parseMoves(tokens *[]string) {
}
