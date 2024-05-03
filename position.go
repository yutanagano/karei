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
	board                  [64]squareState
	colourMasks            [2]bitBoard
	pieceTypeMasks         [6]bitBoard
	kingSquares            [2]coordinate
	enPassantSquare        coordinate
	castlingRights         castlingRights
	activeColour           colour
	pieceColourTypeCounter [12]int
	halfMoveClock          uint8
}

func (p position) occupiedMask() bitBoard {
	return p.colourMasks[white] | p.colourMasks[black]
}

func (p *position) clear() {
	for idx := range p.board {
		p.board[idx] = empty
	}
	for idx := range p.colourMasks {
		p.colourMasks[idx] = 0
	}
	for idx := range p.pieceTypeMasks {
		p.pieceTypeMasks[idx] = 0
	}
	for idx := range p.kingSquares {
		p.kingSquares[idx] = nullCoordinate
	}
	p.enPassantSquare = nullCoordinate
	p.castlingRights = 0
	p.activeColour = white
	for idx := range p.pieceColourTypeCounter {
		p.pieceColourTypeCounter[idx] = 0
	}
	p.halfMoveClock = 0
}

func (p *position) loadFEN(f fen) error {
	p.clear()
	currentRow, currentColumn := 7, 0
	for _, currentRune := range f.boardState {
		if currentColumn > 8 {
			return fmt.Errorf("bad FEN: overfilled row during board specification, row %v col %v", currentRow, currentColumn)
		}

		if currentRune == '/' {
			if currentColumn != 8 {
				return fmt.Errorf("bad FEN: underfilled row during board specification, row %v col %v", currentRow, currentColumn)
			}

			currentRow--
			currentColumn = 0
			continue
		}

		if numEmptySquares, err := strconv.Atoi(string(currentRune)); err == nil {
			currentColumn += numEmptySquares
			continue
		}

		currentCoordinate := coordinateFromRowColumn(currentRow, currentColumn)
		currentSquareState, err := squareStateFromRune(currentRune)

		if err != nil {
			return fmt.Errorf("bad FEN: %s", err.Error())
		}

		p.setSquare(currentSquareState, currentCoordinate)
		currentColumn++
	}

	switch f.activeColour {
	case "w":
		p.activeColour = white
	case "b":
		p.activeColour = black
	default:
		return fmt.Errorf("bad FEN: unrecognised colour %s", f.activeColour)
	}

	p.castlingRights = 0
	switch f.castlingRights {
	case "-":
		break
	default:
		for _, theRune := range f.castlingRights {
			if theRune == 'K' {
				p.castlingRights |= whiteCastleKingSide
				continue
			}
			if theRune == 'Q' {
				p.castlingRights |= whiteCastleQueenSide
				continue
			}
			if theRune == 'k' {
				p.castlingRights |= blackCastleKingSide
				continue
			}
			if theRune == 'q' {
				p.castlingRights |= blackCastleQueenSide
				continue
			}

			return fmt.Errorf("bad FEN: unrecognised character in castling rights specification: %c", theRune)
		}
	}

	switch f.enPassantSquare {
	case "-":
		p.enPassantSquare = nullCoordinate
	default:
		eps, err := coordinateFromString(f.enPassantSquare)
		if err != nil {
			return fmt.Errorf("bad FEN: %s", err.Error())
		}
		p.enPassantSquare = eps
	}

	hmcInt, err := strconv.Atoi(f.halfMoveClock)
	if err != nil {
		return fmt.Errorf("bad FEN: %s", err.Error())
	}
	if hmcInt < 0 {
		return fmt.Errorf("bad FEN: half move clock is negative")
	}
	p.halfMoveClock = uint8(hmcInt)

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
