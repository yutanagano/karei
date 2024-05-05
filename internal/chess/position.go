package chess

import (
	"fmt"
	"strconv"
)

type Position struct {
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

type castlingRights uint8
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

const (
	whiteCastleKingSide  castlingRights = 0b0001
	whiteCastleQueenSide castlingRights = 0b0010
	blackCastleKingSide  castlingRights = 0b0100
	blackCastleQueenSide castlingRights = 0b1000
)

func (p Position) occupiedMask() bitBoard {
	return p.colourMasks[white] | p.colourMasks[black]
}

func (p *Position) clear() {
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

func (p *Position) LoadFEN(f FEN) error {
	p.clear()
	currentRankIndex, currentFileIndex := 7, 0
	for _, currentRune := range f.BoardState {
		if currentFileIndex > 8 {
			return fmt.Errorf("bad FEN: overfilled row during board specification, row %v col %v", currentRankIndex, currentFileIndex)
		}

		if currentRune == '/' {
			if currentFileIndex != 8 {
				return fmt.Errorf("bad FEN: underfilled row during board specification, row %v col %v", currentRankIndex, currentFileIndex)
			}

			currentRankIndex--
			currentFileIndex = 0
			continue
		}

		if numEmptySquares, err := strconv.Atoi(string(currentRune)); err == nil {
			currentFileIndex += numEmptySquares
			continue
		}

		currentCoordinate := coordinateFromRankFileIndices(currentRankIndex, currentFileIndex)
		currentSquareState, err := squareStateFromRune(currentRune)

		if err != nil {
			return fmt.Errorf("bad FEN: %s", err.Error())
		}

		p.setSquare(currentSquareState, currentCoordinate)
		currentFileIndex++
	}

	switch f.ActiveColour {
	case "w":
		p.activeColour = white
	case "b":
		p.activeColour = black
	default:
		return fmt.Errorf("bad FEN: unrecognised colour %s", f.ActiveColour)
	}

	p.castlingRights = 0
	switch f.CastlingRights {
	case "-":
		break
	default:
		for _, theRune := range f.CastlingRights {
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

	switch f.EnPassantSquare {
	case "-":
		p.enPassantSquare = nullCoordinate
	default:
		eps, err := coordinateFromString(f.EnPassantSquare)
		if err != nil {
			return fmt.Errorf("bad FEN: %s", err.Error())
		}
		p.enPassantSquare = eps
	}

	hmcInt, err := strconv.Atoi(f.HalfMoveClock)
	if err != nil {
		return fmt.Errorf("bad FEN: %s", err.Error())
	}
	if hmcInt < 0 {
		return fmt.Errorf("bad FEN: half move clock is negative")
	}
	p.halfMoveClock = uint8(hmcInt)

	return nil
}

func (p *Position) setSquare(state squareState, coord coordinate) {
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

func (p *Position) newGame() {
	// TODO: set up starting position
}

func parseMoves(tokens *[]string) {
}
