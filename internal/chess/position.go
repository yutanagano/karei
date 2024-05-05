package chess

import (
	"fmt"
	"strconv"
)

type Position struct {
	board                  [64]squareState
	occupationByColour     [2]bitBoard
	occupationByPieceType  [6]bitBoard
	controlByColour        [2]bitBoard
	kingSquares            [2]coordinate
	enPassantSquare        coordinate
	castlingRights         castlingRights
	activeColour           colour
	pieceColourTypeCounter [12]int
	halfMoveClock          uint8
}

func (p *Position) LoadFEN(f FEN) error {
	p.clear()
	var currentRankIndex, currentFileIndex int8 = 7, 0
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
			currentFileIndex += int8(numEmptySquares)
			continue
		}

		currentCoordinate, _ := coordinateFromRankFileIndices(currentRankIndex, currentFileIndex)
		currentSquareState, err := squareStateFromRune(currentRune)

		if err != nil {
			return fmt.Errorf("bad FEN: %s", err.Error())
		}

		p.setSquare(currentCoordinate, currentSquareState)
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

func (p *Position) clear() {
	for idx := range p.board {
		p.board[idx] = empty
	}
	for idx := range p.occupationByColour {
		p.occupationByColour[idx] = 0
	}
	for idx := range p.occupationByPieceType {
		p.occupationByPieceType[idx] = 0
	}
	for idx := range p.controlByColour {
		p.controlByColour[idx] = 0
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

func (p *Position) setSquare(theCoord coordinate, theState squareState) {
	p.board[theCoord] = theState

	if theState == empty {
		p.occupationByColour[white].clear(theCoord)
		p.occupationByColour[black].clear(theCoord)

		for thePieceType := king; thePieceType <= pawn; thePieceType++ {
			p.occupationByPieceType[thePieceType].clear(theCoord)
		}
		return
	}

	thePieceType := theState.getPieceType()
	theColour := theState.getColour()

	if thePieceType == king {
		p.kingSquares[theColour] = theCoord
	}

	p.occupationByColour[theColour].set(theCoord)
	p.occupationByPieceType[thePieceType].set(theCoord)
}

func (p Position) getPseudoLegalMoves() moveList {
	pseudoLegalMoves := moveList{}

	pseudoLegalMoves = append(pseudoLegalMoves, p.getPseudoLegalKingMoves()...)
	pseudoLegalMoves = append(pseudoLegalMoves, p.getPseudoLegalQueenMoves()...)
	pseudoLegalMoves = append(pseudoLegalMoves, p.getPseudoLegalRookMoves()...)
	pseudoLegalMoves = append(pseudoLegalMoves, p.getPseudoLegalBishopMoves()...)
	// pseudoLegalMoves = append(pseudoLegalMoves, p.getPseudoLegalKnightMoves()...)
	// pseudoLegalMoves = append(pseudoLegalMoves, p.getPseudoLegalPawnMoves()...)

	return pseudoLegalMoves
}

func (p Position) getPseudoLegalKingMoves() moveList {
	moves := moveList{}

	currentCoord := p.kingSquares[p.activeColour]
	currentRank := currentCoord.getRankIndex()
	currentFile := currentCoord.getFileIndex()

	isLegalDestination := func(toCoord coordinate) bool {
		toSquareState := p.board[toCoord]
		if toSquareState != empty && toSquareState.getColour() == p.activeColour {
			return false
		}
		return true
	}

	for _, rankOffset := range []int8{-1, 0, 1} {
		for _, fileOffset := range []int8{-1, 0, 1} {
			if rankOffset == 0 && fileOffset == 0 {
				continue
			}
			toCoord, err := coordinateFromRankFileIndices(currentRank+rankOffset, currentFile+fileOffset)
			if err != nil {
				continue
			}
			if !isLegalDestination(toCoord) {
				continue
			}
			moves.addMove(currentCoord, toCoord, empty)
		}
	}

	return moves
}

func (p Position) getPseudoLegalQueenMoves() moveList {
	moves := moveList{}

	queensBitBoard := p.occupationByColour[p.activeColour] & p.occupationByPieceType[queen]

	for {
		currentCoord, ok := queensBitBoard.pop()
		if !ok {
			break
		}
		moves = append(moves, p.getPseudoLegalRookMovesFromCoordinate(currentCoord)...)
		moves = append(moves, p.getPseudoLegalBishopMovesFromCoordinate(currentCoord)...)
	}

	return moves
}

func (p Position) getPseudoLegalRookMoves() moveList {
	moves := moveList{}

	rooksBitBoard := p.occupationByColour[p.activeColour] & p.occupationByPieceType[rook]

	for {
		currentCoord, ok := rooksBitBoard.pop()
		if !ok {
			break
		}
		moves = append(moves, p.getPseudoLegalRookMovesFromCoordinate(currentCoord)...)
	}

	return moves
}

func (p Position) getPseudoLegalBishopMoves() moveList {
	moves := moveList{}

	bishopsBitBoard := p.occupationByColour[p.activeColour] & p.occupationByPieceType[bishop]

	for {
		currentCoord, ok := bishopsBitBoard.pop()
		if !ok {
			break
		}
		moves = append(moves, p.getPseudoLegalBishopMovesFromCoordinate(currentCoord)...)
	}

	return moves
}

func (p Position) getPseudoLegalRookMovesFromCoordinate(theCoord coordinate) moveList {
	moves := moveList{}

	currentRank := theCoord.getRankIndex()
	currentFile := theCoord.getFileIndex()

	for toRank := currentRank + 1; toRank < 8; toRank++ {
		toCoord, _ := coordinateFromRankFileIndices(toRank, currentFile)
		deadEnd := p.processLongRangePieceMovesUntilDeadEnd(theCoord, toCoord, &moves)
		if deadEnd {
			break
		}
	}
	for toRank := currentRank - 1; toRank >= 0; toRank-- {
		toCoord, _ := coordinateFromRankFileIndices(toRank, currentFile)
		deadEnd := p.processLongRangePieceMovesUntilDeadEnd(theCoord, toCoord, &moves)
		if deadEnd {
			break
		}
	}
	for toFile := currentFile + 1; toFile < 8; toFile++ {
		toCoord, _ := coordinateFromRankFileIndices(currentRank, toFile)
		deadEnd := p.processLongRangePieceMovesUntilDeadEnd(theCoord, toCoord, &moves)
		if deadEnd {
			break
		}
	}
	for toFile := currentFile - 1; toFile >= 0; toFile-- {
		toCoord, _ := coordinateFromRankFileIndices(currentRank, toFile)
		deadEnd := p.processLongRangePieceMovesUntilDeadEnd(theCoord, toCoord, &moves)
		if deadEnd {
			break
		}
	}

	return moves
}

func (p Position) getPseudoLegalBishopMovesFromCoordinate(theCoord coordinate) moveList {
	moves := moveList{}

	currentRank := theCoord.getRankIndex()
	currentFile := theCoord.getFileIndex()

	for toRank, toFile := currentRank+1, currentFile+1; toRank < 8 && toFile < 8; toRank, toFile = toRank+1, toFile+1 {
		toCoord, _ := coordinateFromRankFileIndices(toRank, toFile)
		deadEnd := p.processLongRangePieceMovesUntilDeadEnd(theCoord, toCoord, &moves)
		if deadEnd {
			break
		}
	}
	for toRank, toFile := currentRank+1, currentFile-1; toRank < 8 && toFile >= 0; toRank, toFile = toRank+1, toFile-1 {
		toCoord, _ := coordinateFromRankFileIndices(toRank, toFile)
		deadEnd := p.processLongRangePieceMovesUntilDeadEnd(theCoord, toCoord, &moves)
		if deadEnd {
			break
		}
	}
	for toRank, toFile := currentRank-1, currentFile+1; toRank >= 0 && toFile < 8; toRank, toFile = toRank-1, toFile+1 {
		toCoord, _ := coordinateFromRankFileIndices(toRank, toFile)
		deadEnd := p.processLongRangePieceMovesUntilDeadEnd(theCoord, toCoord, &moves)
		if deadEnd {
			break
		}
	}
	for toRank, toFile := currentRank-1, currentFile-1; toRank >= 0 && toFile >= 0; toRank, toFile = toRank-1, toFile-1 {
		toCoord, _ := coordinateFromRankFileIndices(toRank, toFile)
		deadEnd := p.processLongRangePieceMovesUntilDeadEnd(theCoord, toCoord, &moves)
		if deadEnd {
			break
		}
	}

	return moves
}

func (p Position) processLongRangePieceMovesUntilDeadEnd(currentCoord, toCoord coordinate, moves *moveList) (deadEnd bool) {
	toSquareState := p.board[toCoord]
	if toSquareState != empty && toSquareState.getColour() == p.activeColour {
		return true
	}
	moves.addMove(currentCoord, toCoord, empty)
	if toSquareState != empty {
		return true
	}
	return false
}

func (p *Position) MakeMove(theMove move) error {
	// check for pseudolegality
	// TODO make this check robust
	fromSquareState := p.board[theMove.From]
	if fromSquareState == empty {
		return fmt.Errorf("no piece to move: %s", theMove.toString())
	}
	if fromSquareState.getColour() != p.activeColour {
		return fmt.Errorf("attempting to move piece of wrong colour: %s", theMove.toString())
	}

	toSquareState := p.board[theMove.To]
	if toSquareState.getColour() == p.activeColour {
		return fmt.Errorf("cannot move piece to square occupied by friendly piece: %s", theMove.toString())
	}

	if theMove.Promotion != empty && theMove.Promotion.getColour() != p.activeColour {
		return fmt.Errorf("cannot promote to enemy piece: %s", theMove.toString())
	}

	// make pseudomove

	return nil
}

func (p *Position) makePseudoMove(theMove move) (resultsInCheck bool) {
	p.enPassantSquare = nullCoordinate

	fromSquareState := p.board[theMove.From]

	switch fromSquareState {
	case whiteKing:
		p.castlingRights.turnOff(whiteCastleKingSide | whiteCastleQueenSide)
		if theMove.From == e1 && theMove.To == c1 {
			p.setSquare(a1, empty)
			p.setSquare(d1, whiteRook)
		} else if theMove.From == e1 && theMove.To == g1 {
			p.setSquare(h1, empty)
			p.setSquare(f1, whiteRook)
		}
	case blackKing:
		p.castlingRights.turnOff(blackCastleKingSide | blackCastleQueenSide)
		if theMove.From == e8 && theMove.To == c8 {
			p.setSquare(a8, empty)
			p.setSquare(d8, whiteRook)
		} else if theMove.From == e8 && theMove.To == g8 {
			p.setSquare(h8, empty)
			p.setSquare(f8, whiteRook)
		}
	case whiteRook:
		if p.castlingRights.isSet(whiteCastleKingSide) && theMove.From == h1 {
			p.castlingRights.turnOff(whiteCastleKingSide)
		} else if p.castlingRights.isSet(whiteCastleQueenSide) && theMove.From == a1 {
			p.castlingRights.turnOff(whiteCastleQueenSide)
		}
	case blackRook:
		if p.castlingRights.isSet(blackCastleKingSide) && theMove.From == h8 {
			p.castlingRights.turnOff(blackCastleKingSide)
		} else if p.castlingRights.isSet(blackCastleQueenSide) && theMove.From == a8 {
			p.castlingRights.turnOff(blackCastleQueenSide)
		}
	case whitePawn:
		switch theMove.getOffset() {
		case 16:
			p.enPassantSquare = theMove.From + 8
		case 7:
			p.setSquare(theMove.From-1, empty)
		case 9:
			p.setSquare(theMove.From+1, empty)
		}
	case blackPawn:
		switch theMove.getOffset() {
		case -16:
			p.enPassantSquare = theMove.From - 8
		case -7:
			p.setSquare(theMove.From+1, empty)
		case -9:
			p.setSquare(theMove.From-1, empty)
		}
	}

	p.setSquare(theMove.From, empty)
	if theMove.Promotion != empty {
		p.setSquare(theMove.To, theMove.Promotion)
	} else {
		p.setSquare(theMove.To, fromSquareState)
	}

	p.activeColour = p.activeColour.getOpposite()

	// TODO check if the player who moved is now in check
	resultsInCheck = false

	return resultsInCheck
}

func (p Position) occupiedMask() bitBoard {
	return p.occupationByColour[white] | p.occupationByColour[black]
}
