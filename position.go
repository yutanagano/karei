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
	board           [64]squareState
	colorMasks      [2]bitBoard
	typeMasks       [6]bitBoard
	kingSquares     [2]coordinate
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
	// TODO: set a square to be some piece
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

func parseMoves(tokens *[]string) {
}
