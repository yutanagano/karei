package chess

import (
	"fmt"
	"strconv"
)

type Position struct {
	board                [64]option[piece]
	bitBoardsByColour    [2]bitBoard
	bitBoardsByPieceType [6]bitBoard
	enPassantSquare      option[coordinate]
	castlingRights       castlingRights
	activeColour         colour
	halfMoveClock        uint8
}

func PositionFromFEN(f FEN) (Position, error) {
	p := Position{}

	currentFile := fileA
	currentRank := rank8
	for _, currentRune := range f.BoardState {
		if currentFile > fileH {
			return p, fmt.Errorf("bad FEN: overfilled row during board specification, row %v col %v", currentRank, currentFile)
		}

		if currentRune == '/' {
			if currentFile != fileH+1 {
				return p, fmt.Errorf("bad FEN: underfilled row during board specification, row %v col %v", currentRank, currentFile)
			}

			currentRank--
			currentFile = fileA
			continue
		}

		if numEmptySquares, err := strconv.Atoi(string(currentRune)); err == nil {
			currentFile += uint8(numEmptySquares)
			continue
		}

		thePiece, err := pieceFromRune(currentRune)
		if err != nil {
			return p, fmt.Errorf("bad FEN: %s", err.Error())
		}

		p.setSquare(coordinateFromParts(currentFile, currentRank), thePiece)
		currentFile++
	}

	switch f.ActiveColour {
	case "w":
		p.activeColour = white
	case "b":
		p.activeColour = black
	default:
		return p, fmt.Errorf("bad FEN: unrecognised colour %s", f.ActiveColour)
	}

	if f.CastlingRights != "-" {
		for _, theRune := range f.CastlingRights {
			if theRune == 'K' {
				p.castlingRights.turnOn(whiteCastleKingSide)
				continue
			}
			if theRune == 'Q' {
				p.castlingRights.turnOn(whiteCastleQueenSide)
				continue
			}
			if theRune == 'k' {
				p.castlingRights.turnOn(blackCastleKingSide)
				continue
			}
			if theRune == 'q' {
				p.castlingRights.turnOn(blackCastleQueenSide)
				continue
			}

			return p, fmt.Errorf("bad FEN: unrecognised character in castling rights specification: %c", theRune)
		}
	}

	if f.EnPassantSquare != "-" {
		eps, err := coordinateFromString(f.EnPassantSquare)
		if err != nil {
			return p, fmt.Errorf("bad FEN: %s", err.Error())
		}
		p.enPassantSquare.setValue(eps)
	}

	hmcInt, err := strconv.Atoi(f.HalfMoveClock)
	if err != nil {
		return p, fmt.Errorf("bad FEN: %s", err.Error())
	}
	if hmcInt < 0 {
		return p, fmt.Errorf("bad FEN: half move clock is negative")
	}
	p.halfMoveClock = uint8(hmcInt)

	return p, nil
}

func (p *Position) getLegalMoves() moveList {
	moves := moveList{}

	enemyAttackBoard := p.getEnemyAttackBoard()
	p.getKingPseudoLegalMoves(&moves)
	p.getQueenPseudoLegalMoves(&moves)
	p.getRookPseudoLegalMoves(&moves)
	p.getBishopPseudoLegalMoves(&moves)
	p.getKnightPseudoLegalMoves(&moves)
	p.getPawnPseudoLegalMoves(&moves)

	moves.filter(p.isLegalMove)

	return moves
}

func (p Position) getEnemyAttackBoard() bitBoard {
	enemyColour := p.activeColour.getOpponent()
	enemyOccupationBoard := p.bitBoardsByColour[enemyColour]
	generalOccupationBoard := p.getOccupationBitBoard()
	enemyAttackBoard := bitBoard(0)

	enemyKingBoard := enemyOccupationBoard & p.bitBoardsByPieceType[king]
	enemyKingCoord, _ := enemyKingBoard.pop()
	enemyAttackBoard |= getKingControlBitBoard(enemyKingCoord)

	enemyQueenBoard := enemyOccupationBoard & p.bitBoardsByPieceType[queen]
	for {
		enemyQueenCoord, ok := enemyQueenBoard.pop()
		if !ok {
			break
		}
		enemyAttackBoard |= getQueenControlBitBoard(enemyQueenCoord, generalOccupationBoard)
	}

	enemyRookBoard := enemyOccupationBoard & p.bitBoardsByPieceType[rook]
	for {
		enemyRookCoord, ok := enemyRookBoard.pop()
		if !ok {
			break
		}
		enemyAttackBoard |= getRookControlBitBoard(enemyRookCoord, generalOccupationBoard)
	}

	enemyBishopBoard := enemyOccupationBoard & p.bitBoardsByPieceType[bishop]
	for {
		enemyBishopCoord, ok := enemyBishopBoard.pop()
		if !ok {
			break
		}
		enemyAttackBoard |= getBishopControlBitBoard(enemyBishopCoord, generalOccupationBoard)
	}

	enemyKnightBoard := enemyOccupationBoard & p.bitBoardsByPieceType[knight]
	for {
		enemyKnightCoord, ok := enemyKnightBoard.pop()
		if !ok {
			break
		}
		enemyAttackBoard |= getKnightControlBitBoard(enemyKnightCoord)
	}

	enemyPawnBoard := enemyOccupationBoard & p.bitBoardsByPieceType[pawn]
	enemyAttackBoard |= getPawnKingSideControlBitBoard(enemyPawnBoard, enemyColour)
	enemyAttackBoard |= getPawnQueenSideControlBitBoard(enemyPawnBoard, enemyColour)

	return enemyAttackBoard
}

func (p Position) getKingPseudoLegalMoves(enemyAttackBoard bitBoard, moveCandidates *moveList) {
	kingBitBoard := p.bitBoardsByColour[p.activeColour] & p.bitBoardsByPieceType[king]
	currentCoord, _ := kingBitBoard.pop()
	controlledSquares := getKingControlBitBoard(currentCoord)

	for {
		toCoord, ok := controlledSquares.pop()
		if !ok {
			break
		}

		if !enemyAttackBoard.get(toCoord) && !p.isOccupiedByFriendly(toCoord) {
			newMove := p.moveFromAlgebraicParts(currentCoord, toCoord, newEmptyOption[piece]())
			moveCandidates.add(newMove)
		}
	}

	if enemyAttackBoard.get(currentCoord) {
		return
	}

	clearForCastling := func(c coordinate) bool {
		return !(p.getOccupationBitBoard() | enemyAttackBoard).get(c)
	}

	switch p.activeColour {
	case white:
		if p.castlingRights.isSet(whiteCastleKingSide) && clearForCastling(f1) && clearForCastling(g1) {
			whiteKingSideCastle := p.moveFromAlgebraicParts(e1, g1, newEmptyOption[piece]())
			moveCandidates.add(whiteKingSideCastle)
		}
		if p.castlingRights.isSet(whiteCastleQueenSide) && clearForCastling(d1) && clearForCastling(c1) {
			whiteQueenSideCastle := p.moveFromAlgebraicParts(e1, g1, newEmptyOption[piece]())
			moveCandidates.add(whiteQueenSideCastle)
		}
	case black:
		if p.castlingRights.isSet(blackCastleKingSide) && clearForCastling(f8) && clearForCastling(g8) {
			blackKingSideCastle := p.moveFromAlgebraicParts(e1, g1, newEmptyOption[piece]())
			moveCandidates.add(blackKingSideCastle)
		}
		if p.castlingRights.isSet(blackCastleQueenSide) && clearForCastling(d8) && clearForCastling(c8) {
			blackQueenSideCastle := p.moveFromAlgebraicParts(e1, g1, newEmptyOption[piece]())
			moveCandidates.add(blackQueenSideCastle)
		}
	}
}

func (p Position) getQueenPseudoLegalMoves(moveCandidates *moveList) {
	occupiedSquares := p.getOccupationBitBoard()

	queensBitBoard := p.bitBoardsByColour[p.activeColour] & p.bitBoardsByPieceType[queen]
	for {
		currentCoord, ok := queensBitBoard.pop()
		if !ok {
			break
		}

		controlBitBoard := getQueenControlBitBoard(currentCoord, occupiedSquares)
		for {
			toCoord, ok := controlBitBoard.pop()
			if !ok {
				break
			}
			if !p.isOccupiedByFriendly(toCoord) {
				newMove := p.moveFromAlgebraicParts(currentCoord, toCoord, newEmptyOption[piece]())
				moveCandidates.add(newMove)
			}
		}
	}
}

func (p Position) getRookPseudoLegalMoves(moveCandidates *moveList) {
	occupiedSquares := p.getOccupationBitBoard()

	rooksBitBoard := p.bitBoardsByColour[p.activeColour] & p.bitBoardsByPieceType[rook]
	for {
		currentCoord, ok := rooksBitBoard.pop()
		if !ok {
			break
		}

		controlBitBoard := getRookControlBitBoard(currentCoord, occupiedSquares)
		for {
			toCoord, ok := controlBitBoard.pop()
			if !ok {
				break
			}
			if !p.isOccupiedByFriendly(toCoord) {
				newMove := p.moveFromAlgebraicParts(currentCoord, toCoord, newEmptyOption[piece]())
				moveCandidates.add(newMove)
			}
		}
	}
}

func (p Position) getBishopPseudoLegalMoves(moveCandidates *moveList) {
	occupiedSquares := p.getOccupationBitBoard()

	bishopsBitBoard := p.bitBoardsByColour[p.activeColour] & p.bitBoardsByPieceType[bishop]
	for {
		currentCoord, ok := bishopsBitBoard.pop()
		if !ok {
			break
		}

		controlBitBoard := getBishopControlBitBoard(currentCoord, occupiedSquares)
		for {
			toCoord, ok := controlBitBoard.pop()
			if !ok {
				break
			}
			if !p.isOccupiedByFriendly(toCoord) {
				newMove := p.moveFromAlgebraicParts(currentCoord, toCoord, newEmptyOption[piece]())
				moveCandidates.add(newMove)
			}
		}
	}
}

func (p Position) getKnightPseudoLegalMoves(moveCandidates *moveList) {
	knightsBitBoard := p.bitBoardsByColour[p.activeColour] & p.bitBoardsByPieceType[knight]
	for {
		currentCoord, ok := knightsBitBoard.pop()
		if !ok {
			break
		}

		controlBitBoard := getKnightControlBitBoard(currentCoord)
		for {
			toCoord, ok := controlBitBoard.pop()
			if !ok {
				break
			}
			if !p.isOccupiedByFriendly(toCoord) {
				newMove := p.moveFromAlgebraicParts(currentCoord, toCoord, newEmptyOption[piece]())
				moveCandidates.add(newMove)
			}
		}
	}
}

func (p *Position) surveyBishopActivity(c colour, moveCandidates *moveList, enemyAttackBoard *bitBoard) {
	bishopsBitBoard := p.bitBoardsByColour[c] & p.bitBoardsByPieceType[bishop]
	for {
		currentCoord, ok := bishopsBitBoard.pop()
		if !ok {
			break
		}

		for _, o := range []gridOffset{
			{1, 1},
			{-1, 1},
			{-1, -1},
			{1, -1},
		} {
			p.surveySlidingControlFromCoordinate(currentCoord, o, c, moveCandidates, enemyAttackBoard)
		}
	}
}

func (p *Position) surveyKnightActivity(c colour, moveCandidates *moveList, enemyAttackBoard *bitBoard) {
	isForEnemy := c != p.activeColour
	knightsBitBoard := p.bitBoardsByColour[c] & p.bitBoardsByPieceType[knight]

	for {
		currentCoord, ok := knightsBitBoard.pop()
		if !ok {
			break
		}

		controlledSquares := knightControlFrom[currentCoord]

		if isForEnemy {
			*enemyAttackBoard |= controlledSquares
			continue
		}

		for {
			toCoord, ok := controlledSquares.pop()
			if !ok {
				break
			}

			if !p.isOccupiedByFriendly(c, toCoord) {
				newMove := p.moveFromAlgebraicParts(currentCoord, toCoord, empty)
				moveCandidates.add(newMove)
			}
		}
	}
}

func (p *Position) surveyPawnActivity(c colour, moveCandidates *moveList, enemyAttackBoard *bitBoard) {
	switch c {
	case white:
		p.surveyPawnActivityWhite(moveCandidates, enemyAttackBoard)
	case black:
		p.surveyPawnActivityBlack(moveCandidates, enemyAttackBoard)
	}
}

func (p *Position) surveyPawnActivityWhite(moveCandidates *moveList, enemyAttackBoard *bitBoard) {
	isForEnemy := p.activeColour == black
	pawnBitBoard := p.bitBoardsByColour[white] & p.bitBoardsByPieceType[pawn]

	kingSideControl := (pawnBitBoard & ^bitBoardFileHMask) << 9
	queenSideControl := (pawnBitBoard & ^bitBoardFileAMask) << 7

	if isForEnemy {
		*enemyAttackBoard |= kingSideControl | queenSideControl
		return
	}

	capturableSquares := p.bitBoardsByColour[black]
	if !p.enPassantSquare.isEmpty() {
		capturableSquares.turnOn(p.enPassantSquare.getValue())
	}

	occupiedSquares := p.getOccupationBitBoard()
	addWhitePawnMoves := func(from coordinate, to coordinate) {
		if to.getRank() != 7 {
			moveCandidates.add(p.moveFromAlgebraicParts(from, to, empty))
			return
		}

		for _, promotionPiece := range []squareState{whiteQueen, whiteRook, whiteBishop, whiteKnight} {
			promotion := p.moveFromAlgebraicParts(from, to, promotionPiece)
			moveCandidates.add(promotion)
		}
	}

	kingSideCaptures := kingSideControl & capturableSquares
	for {
		toCoord, ok := kingSideCaptures.pop()
		if !ok {
			break
		}

		fromCoord := toCoord - 9
		addWhitePawnMoves(fromCoord, toCoord)
	}

	queenSideCaptures := queenSideControl & capturableSquares
	for {
		toCoord, ok := queenSideCaptures.pop()
		if !ok {
			break
		}

		fromCoord := toCoord - 7
		addWhitePawnMoves(fromCoord, toCoord)
	}

	oneSquareForward := (pawnBitBoard << 8) & ^occupiedSquares
	twoSquaresForward := (oneSquareForward << 8) & bitBoardRank4Mask & ^occupiedSquares

	for {
		toCoord, ok := oneSquareForward.pop()
		if !ok {
			break
		}

		fromCoord := toCoord - 8
		addWhitePawnMoves(fromCoord, toCoord)
	}

	for {
		toCoord, ok := twoSquaresForward.pop()
		if !ok {
			break
		}

		fromCoord := toCoord - 16
		newMove := p.moveFromAlgebraicParts(fromCoord, toCoord, empty)
		moveCandidates.add(newMove)
	}
}

func (p *Position) surveyPawnActivityBlack(moveCandidates *moveList, enemyAttackBoard *bitBoard) {
	isForEnemy := p.activeColour == white
	pawnBitBoard := p.bitBoardsByColour[black] & p.bitBoardsByPieceType[pawn]

	kingSideControl := (pawnBitBoard & ^bitBoardFileHMask) >> 7
	queenSideControl := (pawnBitBoard & ^bitBoardFileAMask) >> 9

	if isForEnemy {
		*enemyAttackBoard |= kingSideControl | queenSideControl
		return
	}

	capturableSquares := p.bitBoardsByColour[white]
	if !p.enPassantSquare.isEmpty() {
		capturableSquares.turnOn(p.enPassantSquare.getValue())
	}

	occupiedSquares := p.getOccupationBitBoard()
	addBlackPawnMoves := func(from coordinate, to coordinate) {
		if to.getRank() != 0 {
			newMove := p.moveFromAlgebraicParts(from, to, empty)
			moveCandidates.add(newMove)
			return
		}

		for _, promotionPiece := range []squareState{blackQueen, blackRook, blackBishop, blackKnight} {
			promotion := p.moveFromAlgebraicParts(from, to, promotionPiece)
			moveCandidates.add(promotion)
		}
	}

	kingSideCaptures := kingSideControl & capturableSquares
	for {
		toCoord, ok := kingSideCaptures.pop()
		if !ok {
			break
		}

		fromCoord := toCoord + 7
		addBlackPawnMoves(fromCoord, toCoord)
	}

	queenSideCaptures := queenSideControl & capturableSquares
	for {
		toCoord, ok := queenSideCaptures.pop()
		if !ok {
			break
		}

		fromCoord := toCoord + 9
		addBlackPawnMoves(fromCoord, toCoord)
	}

	oneSquareForward := (pawnBitBoard >> 8) & ^occupiedSquares
	twoSquaresForward := (oneSquareForward >> 8) & bitBoardRank5Mask & ^occupiedSquares

	for {
		toCoord, ok := oneSquareForward.pop()
		if !ok {
			break
		}

		fromCoord := toCoord + 8
		addBlackPawnMoves(fromCoord, toCoord)
	}

	for {
		toCoord, ok := twoSquaresForward.pop()
		if !ok {
			break
		}

		fromCoord := toCoord + 16
		newMove := p.moveFromAlgebraicParts(fromCoord, toCoord, empty)
		moveCandidates.add(newMove)
	}
}

func (p Position) moveFromAlgebraicParts(from, to coordinate, promotionTo option[piece]) move {
	return moveFromParts(from, to, p.getSquare(to).getValue(), promotionTo, p.castlingRights, p.enPassantSquare)
}

func (p *Position) isLegalMove(theMove move) bool {
	return true
}

func (p *Position) makePseudoLegalMove(theMove move) {
	p.enPassantSquare = nullCoordinate

	fromCoord := theMove.getFromCoordinate()
	toCoord := theMove.getToCoordinate()
	EPSquare := theMove.getCurrentEPSquare()
	pieceBeingMoved := p.getSquare(fromCoord).getValue()

	switch pieceBeingMoved {
	case whiteKing:
		p.castlingRights.turnOff(whiteCastleKingSide | whiteCastleQueenSide)
		if fromCoord == e1 && toCoord == c1 {
			p.setSquare(a1, empty)
			p.setSquare(d1, whiteRook)
		} else if fromCoord == e1 && toCoord == g1 {
			p.setSquare(h1, empty)
			p.setSquare(f1, whiteRook)
		}
	case blackKing:
		p.castlingRights.turnOff(blackCastleKingSide | blackCastleQueenSide)
		if fromCoord == e8 && toCoord == c8 {
			p.setSquare(a8, empty)
			p.setSquare(d8, blackRook)
		} else if fromCoord == e8 && toCoord == g8 {
			p.setSquare(h8, empty)
			p.setSquare(f8, blackRook)
		}
	case whiteRook:
		if p.castlingRights.isSet(whiteCastleKingSide) && fromCoord == h1 {
			p.castlingRights.turnOff(whiteCastleKingSide)
		} else if p.castlingRights.isSet(whiteCastleQueenSide) && fromCoord == a1 {
			p.castlingRights.turnOff(whiteCastleQueenSide)
		}
	case blackRook:
		if p.castlingRights.isSet(blackCastleKingSide) && fromCoord == h8 {
			p.castlingRights.turnOff(blackCastleKingSide)
		} else if p.castlingRights.isSet(blackCastleQueenSide) && fromCoord == a8 {
			p.castlingRights.turnOff(blackCastleQueenSide)
		}
	case whitePawn:
		if toCoord-fromCoord == 16 {
			p.enPassantSquare = fromCoord + 8
		}
		if toCoord == EPSquare {
			p.setSquare(toCoord-8, empty)
		}
	case blackPawn:
		if fromCoord-toCoord == 16 {
			p.enPassantSquare = fromCoord - 8
		}
		if toCoord == EPSquare {
			p.setSquare(toCoord+8, empty)
		}
	}

	p.setSquare(fromCoord, empty)
	if promotionTo := theMove.getPromotionTo(); promotionTo != empty {
		p.setSquare(toCoord, promotionTo)
	} else {
		p.setSquare(toCoord, pieceBeingMoved)
	}

	p.activeColour = p.activeColour.getOpponent()
}

func (p *Position) unmakeMove(theMove move) {
	p.enPassantSquare = theMove.getCurrentEPSquare()
	p.castlingRights = theMove.getCurrentCastlingRights()

	fromCoord := theMove.getFromCoordinate()
	toCoord := theMove.getToCoordinate()
	pieceBeingMoved := p.getSquare(toCoord)

	switch pieceBeingMoved {
	case whiteKing:
		if fromCoord == e1 && toCoord == c1 {
			p.setSquare(a1, whiteRook)
			p.setSquare(d1, empty)
		} else if fromCoord == e1 && toCoord == g1 {
			p.setSquare(h1, whiteRook)
			p.setSquare(f1, empty)
		}
	case blackKing:
		if fromCoord == e8 && toCoord == c8 {
			p.setSquare(a8, blackRook)
			p.setSquare(d8, empty)
		} else if fromCoord == e8 && toCoord == g8 {
			p.setSquare(h8, blackRook)
			p.setSquare(f8, empty)
		}
	}
}

func (p Position) isOccupiedByFriendly(c coordinate) bool {
	return p.bitBoardsByColour[p.activeColour].get(c)
}

func (p Position) getOccupationBitBoard() bitBoard {
	return p.bitBoardsByColour[white] | p.bitBoardsByColour[black]
}

func (p Position) getSquare(c coordinate) option[piece] {
	return p.board[c]
}

func (p *Position) clearSquare(c coordinate) {
	p.board[c].clear()
	p.bitBoardsByColour[white].turnOff(c)
	p.bitBoardsByColour[black].turnOff(c)

	for idx := range p.bitBoardsByPieceType {
		p.bitBoardsByPieceType[idx].turnOff(c)
	}
}

func (p *Position) setSquare(c coordinate, pc piece) {
	p.board[c].setValue(pc)
	p.bitBoardsByColour[pc.getColour()].turnOn(c)
	p.bitBoardsByPieceType[pc.getPieceType()].turnOn(c)
}

func (p *Position) MakeMove(theMove algebraicMove) error {
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
