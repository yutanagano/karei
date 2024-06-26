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
	legalMoves             moveList
}

func (p *Position) LoadFEN(f FEN) error {
	p.clear()
	var currentFileIndex, currentRankIndex int8 = 0, 7
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

		currentCoordinate, _ := coordinateFromRankFileIndices(currentFileIndex, currentRankIndex)
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

	p.doStaticAnalysis()

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
	p.legalMoves = moveList{}
}

func (p *Position) setSquare(theCoord coordinate, theState squareState) {
	p.board[theCoord] = theState

	if theState == empty {
		p.occupationByColour[white].turnOff(theCoord)
		p.occupationByColour[black].turnOff(theCoord)

		for thePieceType := king; thePieceType <= pawn; thePieceType++ {
			p.occupationByPieceType[thePieceType].turnOff(theCoord)
		}
		return
	}

	thePieceType := theState.getPieceType()
	theColour := theState.getColour()

	if thePieceType == king {
		p.kingSquares[theColour] = theCoord
	}

	p.occupationByColour[theColour].turnOn(theCoord)
	p.occupationByPieceType[thePieceType].turnOn(theCoord)
}

func (p Position) getSquare(theCoord coordinate) squareState {
	return p.board[theCoord]
}

func (p *Position) doStaticAnalysis() {
	p.controlByColour[white] = 0
	p.controlByColour[black] = 0

	p.surveyPieceActivity(p.activeColour.getOpponent(), false)
	p.legalMoves = p.surveyPieceActivity(p.activeColour, true)
	p.legalMoves.filter(p.isLegalMove)
}

func (p *Position) surveyPieceActivity(player colour, getPsuedoLegalMoves bool) moveList {
	pseudoLegalMoves := moveList{}

	p.surveyKingActivity(player, getPsuedoLegalMoves, &pseudoLegalMoves)
	p.surveyQueenActivity(player, getPsuedoLegalMoves, &pseudoLegalMoves)
	p.surveyRookActivity(player, getPsuedoLegalMoves, &pseudoLegalMoves)
	p.surveyBishopActivity(player, getPsuedoLegalMoves, &pseudoLegalMoves)
	p.surveyKnightActivity(player, getPsuedoLegalMoves, &pseudoLegalMoves)
	p.surveyPawnActivity(player, getPsuedoLegalMoves, &pseudoLegalMoves)

	return pseudoLegalMoves
}

func (p *Position) surveyKingActivity(player colour, getPsuedoLegalMoves bool, pseudoLegalMoves *moveList) {
	currentCoord := p.kingSquares[player]
	controlledSquares := kingControlFrom[currentCoord]
	p.controlByColour[player] |= controlledSquares

	if !getPsuedoLegalMoves {
		return
	}

	for {
		toCoord, ok := controlledSquares.pop()
		if !ok {
			break
		}

		if !p.isAttackedByEnemy(player, toCoord) && !p.isOccupiedByFriendly(player, toCoord) {
			newMove := p.moveFromAlgebraicParts(currentCoord, toCoord, empty)
			pseudoLegalMoves.add(newMove)
		}
	}

	if p.inCheck(player) {
		return
	}

	switch player {
	case white:
		if p.castlingRights.isSet(whiteCastleKingSide) && p.allowsSafePassage(white, f1) && p.allowsSafePassage(white, g1) {
			whiteKingSideCastle := p.moveFromAlgebraicParts(e1, g1, empty)
			pseudoLegalMoves.add(whiteKingSideCastle)
		}
		if p.castlingRights.isSet(whiteCastleQueenSide) && p.allowsSafePassage(white, d1) && p.allowsSafePassage(white, c1) {
			whiteQueenSideCastle := p.moveFromAlgebraicParts(e1, g1, empty)
			pseudoLegalMoves.add(whiteQueenSideCastle)
		}
	case black:
		if p.castlingRights.isSet(blackCastleKingSide) && p.allowsSafePassage(black, f8) && p.allowsSafePassage(black, g8) {
			blackKingSideCastle := p.moveFromAlgebraicParts(e1, g1, empty)
			pseudoLegalMoves.add(blackKingSideCastle)
		}
		if p.castlingRights.isSet(blackCastleQueenSide) && p.allowsSafePassage(black, d8) && p.allowsSafePassage(black, c8) {
			blackQueenSideCastle := p.moveFromAlgebraicParts(e1, g1, empty)
			pseudoLegalMoves.add(blackQueenSideCastle)
		}
	}
}

func (p *Position) surveyQueenActivity(player colour, getPsuedoLegalMoves bool, pseudoLegalMoves *moveList) {
	queensBitBoard := p.occupationByColour[player] & p.occupationByPieceType[queen]
	for {
		currentCoord, ok := queensBitBoard.pop()
		if !ok {
			break
		}

		for _, delta := range []gridDelta{
			{1, 0},
			{1, 1},
			{0, 1},
			{-1, 1},
			{-1, 0},
			{-1, -1},
			{0, -1},
			{1, -1},
		} {
			p.surveySlidingControlFromCoordinate(currentCoord, delta, player, getPsuedoLegalMoves, pseudoLegalMoves)
		}
	}
}

func (p *Position) surveyRookActivity(player colour, getPsuedoLegalMoves bool, pseudoLegalMoves *moveList) {
	rooksBitBoard := p.occupationByColour[player] & p.occupationByPieceType[rook]
	for {
		currentCoord, ok := rooksBitBoard.pop()
		if !ok {
			break
		}

		for _, theOffset := range []gridDelta{
			{1, 0},
			{0, 1},
			{-1, 0},
			{0, -1},
		} {
			p.surveySlidingControlFromCoordinate(currentCoord, theOffset, player, getPsuedoLegalMoves, pseudoLegalMoves)
		}
	}
}

func (p *Position) surveyBishopActivity(player colour, getPsuedoLegalMoves bool, pseudoLegalMoves *moveList) {
	bishopsBitBoard := p.occupationByColour[player] & p.occupationByPieceType[bishop]
	for {
		currentCoord, ok := bishopsBitBoard.pop()
		if !ok {
			break
		}

		for _, theOffset := range []gridDelta{
			{1, 1},
			{-1, 1},
			{-1, -1},
			{1, -1},
		} {
			p.surveySlidingControlFromCoordinate(currentCoord, theOffset, player, getPsuedoLegalMoves, pseudoLegalMoves)
		}
	}
}

func (p *Position) surveySlidingControlFromCoordinate(originalCoord coordinate, unitDelta gridDelta, player colour, getPsuedoLegalMoves bool, pseudoLegalMoves *moveList) {
	for toCoord, err := originalCoord.move(unitDelta); err == nil; toCoord, err = toCoord.move(unitDelta) {
		p.controlByColour[player].turnOn(toCoord)
		if getPsuedoLegalMoves && !p.isOccupiedByFriendly(player, toCoord) {
			newMove := p.moveFromAlgebraicParts(originalCoord, toCoord, empty)
			pseudoLegalMoves.add(newMove)
		}
		if p.getSquare(toCoord) != empty {
			break
		}
	}
}

func (p *Position) surveyKnightActivity(player colour, getPsuedoLegalMoves bool, pseudoLegalMoves *moveList) {
	knightsBitBoard := p.occupationByColour[player] & p.occupationByPieceType[knight]

	for {
		currentCoord, ok := knightsBitBoard.pop()
		if !ok {
			break
		}

		controlledSquares := knightControlFrom[currentCoord]
		p.controlByColour[player] |= controlledSquares

		if !getPsuedoLegalMoves {
			return
		}

		for {
			toCoord, ok := controlledSquares.pop()
			if !ok {
				break
			}

			if !p.isOccupiedByFriendly(player, toCoord) {
				newMove := p.moveFromAlgebraicParts(currentCoord, toCoord, empty)
				pseudoLegalMoves.add(newMove)
			}
		}
	}
}

func (p *Position) surveyPawnActivity(player colour, getPsuedoLegalMoves bool, pseudoLegalMoves *moveList) {
	switch player {
	case white:
		p.surveyPawnActivityWhite(getPsuedoLegalMoves, pseudoLegalMoves)
	case black:
		p.surveyPawnActivityBlack(getPsuedoLegalMoves, pseudoLegalMoves)
	}
}

func (p *Position) surveyPawnActivityWhite(getPsuedoLegalMoves bool, pseudoLegalMoves *moveList) {
	pawnBitBoard := p.occupationByColour[white] & p.occupationByPieceType[pawn]

	kingSideControl := (pawnBitBoard & ^fileH) << 9
	queenSideControl := (pawnBitBoard & ^fileA) << 7
	p.controlByColour[white] |= kingSideControl | queenSideControl

	if !getPsuedoLegalMoves {
		return
	}

	capturableSquares := p.occupationByColour[black]
	if p.enPassantSquare != nullCoordinate {
		capturableSquares.turnOn(p.enPassantSquare)
	}

	occupiedSquares := p.getOccupationBitBoard()
	addWhitePawnMoves := func(from coordinate, to coordinate) {
		if to.getRankIndex() != 7 {
			pseudoLegalMoves.add(p.moveFromAlgebraicParts(from, to, empty))
			return
		}

		for _, promotionPiece := range []squareState{whiteQueen, whiteRook, whiteBishop, whiteKnight} {
			promotion := p.moveFromAlgebraicParts(from, to, promotionPiece)
			pseudoLegalMoves.add(promotion)
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
	twoSquaresForward := (oneSquareForward << 8) & rank4 & ^occupiedSquares

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
		pseudoLegalMoves.add(newMove)
	}
}

func (p *Position) surveyPawnActivityBlack(getPsuedoLegalMoves bool, pseudoLegalMoves *moveList) {
	pawnBitBoard := p.occupationByColour[black] & p.occupationByPieceType[pawn]

	kingSideControl := (pawnBitBoard & ^fileH) >> 7
	queenSideControl := (pawnBitBoard & ^fileA) >> 9
	p.controlByColour[black] |= kingSideControl | queenSideControl

	if !getPsuedoLegalMoves {
		return
	}

	capturableSquares := p.occupationByColour[white]
	if p.enPassantSquare != nullCoordinate {
		capturableSquares.turnOn(p.enPassantSquare)
	}

	occupiedSquares := p.getOccupationBitBoard()
	addBlackPawnMoves := func(from coordinate, to coordinate) {
		if to.getRankIndex() != 0 {
			newMove := p.moveFromAlgebraicParts(from, to, empty)
			pseudoLegalMoves.add(newMove)
			return
		}

		for _, promotionPiece := range []squareState{blackQueen, blackRook, blackBishop, blackKnight} {
			promotion := p.moveFromAlgebraicParts(from, to, promotionPiece)
			pseudoLegalMoves.add(promotion)
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
	twoSquaresForward := (oneSquareForward >> 8) & rank5 & ^occupiedSquares

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
		pseudoLegalMoves.add(newMove)
	}
}

func (p Position) moveFromAlgebraicParts(from, to coordinate, promotionTo squareState) move {
	return moveFromParts(from, to, p.getSquare(to), promotionTo, p.castlingRights, p.enPassantSquare)
}

func (p *Position) isLegalMove(theMove move) bool {
	return true
}

func (p *Position) makePseudoLegalMove(theMove move) {
	p.enPassantSquare = nullCoordinate

	fromCoord := theMove.getFromCoordinate()
	toCoord := theMove.getToCoordinate()
	EPSquare := theMove.getCurrentEPSquare()
	pieceBeingMoved := p.getSquare(fromCoord)

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

	p.doStaticAnalysis()
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

func (p Position) inCheck(player colour) bool {
	return p.isAttackedByEnemy(player, p.kingSquares[player])
}

func (p Position) isAttackedByEnemy(player colour, theCoord coordinate) bool {
	return p.controlByColour[player.getOpponent()].get(theCoord)
}

func (p Position) isOccupiedByFriendly(player colour, theCoord coordinate) bool {
	return p.occupationByColour[player].get(theCoord)
}

func (p Position) allowsSafePassage(player colour, theCoord coordinate) bool {
	return !p.getOccupationBitBoard().get(theCoord) && !p.isAttackedByEnemy(player, theCoord)
}

func (p Position) getOccupationBitBoard() bitBoard {
	return p.occupationByColour[white] | p.occupationByColour[black]
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
