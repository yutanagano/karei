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

func (p *Position) doStaticAnalysis() {
	var legalMoves moveList

	p.surveyControlledSquares(p.activeColour.getOpponent(), false)
	legalMoves = p.surveyControlledSquares(p.activeColour, true)

	p.legalMoves = legalMoves
}

func (p *Position) surveyControlledSquares(player colour, getMoveCandidates bool) moveList {
	moveCandidates := moveList{}

	p.surveyKingControl(player, getMoveCandidates, &moveCandidates)
	p.surveyQueenControl(player, getMoveCandidates, &moveCandidates)
	p.surveyRookControl(player, getMoveCandidates, &moveCandidates)
	p.surveyBishopControl(player, getMoveCandidates, &moveCandidates)
	p.surveyKnightControl(player, getMoveCandidates, &moveCandidates)
	p.surveyPawnControl(player, getMoveCandidates, &moveCandidates)

	return moveCandidates
}

func (p *Position) surveyKingControl(player colour, getMoveCandidates bool, moveCandidates *moveList) {
	currentCoord := p.kingSquares[player]
	controlledSquares := kingControlFrom[currentCoord]
	p.controlByColour[player] |= controlledSquares

	if !getMoveCandidates {
		return
	}

	for {
		toCoord, ok := controlledSquares.pop()
		if !ok {
			break
		}

		if !p.isAttackedByEnemy(player, toCoord) && !p.isOccupiedByFriendly(player, toCoord) {
			moveCandidates.addMove(currentCoord, toCoord, empty)
		}
	}
}

func (p *Position) surveyQueenControl(player colour, getMoveCandidates bool, moveCandidates *moveList) {
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
			p.surveySlidingControlFromCoordinate(currentCoord, delta, player, getMoveCandidates, moveCandidates)
		}
	}
}

func (p *Position) surveyRookControl(player colour, getMoveCandidates bool, moveCandidates *moveList) {
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
			p.surveySlidingControlFromCoordinate(currentCoord, theOffset, player, getMoveCandidates, moveCandidates)
		}
	}
}

func (p *Position) surveyBishopControl(player colour, getMoveCandidates bool, moveCandidates *moveList) {
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
			p.surveySlidingControlFromCoordinate(currentCoord, theOffset, player, getMoveCandidates, moveCandidates)
		}
	}
}

func (p *Position) surveySlidingControlFromCoordinate(originalCoord coordinate, unitDelta gridDelta, player colour, getMoveCandidates bool, moveCandidates *moveList) {
	for toCoord, err := originalCoord.move(unitDelta); err == nil; toCoord, err = toCoord.move(unitDelta) {
		p.controlByColour[player].turnOn(toCoord)
		if getMoveCandidates && !p.isOccupiedByFriendly(player, toCoord) {
			moveCandidates.addMove(originalCoord, toCoord, empty)
		}
		if p.board[toCoord] != empty {
			break
		}
	}
}

func (p *Position) surveyKnightControl(player colour, getMoveCandidates bool, moveCandidates *moveList) {
	knightsBitBoard := p.occupationByColour[player] & p.occupationByPieceType[knight]

	for {
		currentCoord, ok := knightsBitBoard.pop()
		if !ok {
			break
		}

		controlledSquares := knightControlFrom[currentCoord]
		p.controlByColour[player] |= controlledSquares

		if !getMoveCandidates {
			return
		}

		for {
			toCoord, ok := controlledSquares.pop()
			if !ok {
				break
			}

			if !p.isOccupiedByFriendly(player, toCoord) {
				moveCandidates.addMove(currentCoord, toCoord, empty)
			}
		}
	}
}

func (p *Position) surveyPawnControl(player colour, getMoveCandidates bool, moveCandidates *moveList) {
	switch player {
	case white:
		p.surveyPawnControlWhite(getMoveCandidates, moveCandidates)
	case black:
		p.surveyPawnControlBlack(getMoveCandidates, moveCandidates)
	}
}

func (p *Position) surveyPawnControlWhite(getMoveCandidates bool, moveCandidates *moveList) {
	pawnBitBoard := p.occupationByColour[white] & p.occupationByPieceType[pawn]

	kingSideControl := (pawnBitBoard & ^fileH) << 9
	queenSideControl := (pawnBitBoard & ^fileA) << 7
	p.controlByColour[white] |= kingSideControl | queenSideControl

	if !getMoveCandidates {
		return
	}

	capturableSquares := p.occupationByColour[black]
	if p.enPassantSquare != nullCoordinate {
		capturableSquares.turnOn(p.enPassantSquare)
	}

	occupiedSquares := p.getOccupationBitBoard()
	addWhitePawnMoves := func(from coordinate, to coordinate) {
		if to.getRankIndex() != 7 {
			moveCandidates.addMove(from, to, empty)
			return
		}

		moveCandidates.addMove(from, to, whiteQueen)
		moveCandidates.addMove(from, to, whiteRook)
		moveCandidates.addMove(from, to, whiteBishop)
		moveCandidates.addMove(from, to, whiteKnight)
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
		moveCandidates.addMove(fromCoord, toCoord, empty)
	}
}

func (p *Position) surveyPawnControlBlack(getMoveCandidates bool, moveCandidates *moveList) {
	pawnBitBoard := p.occupationByColour[black] & p.occupationByPieceType[pawn]

	kingSideControl := (pawnBitBoard & ^fileH) >> 7
	queenSideControl := (pawnBitBoard & ^fileA) >> 9
	p.controlByColour[black] |= kingSideControl | queenSideControl

	if !getMoveCandidates {
		return
	}

	capturableSquares := p.occupationByColour[white]
	if p.enPassantSquare != nullCoordinate {
		capturableSquares.turnOn(p.enPassantSquare)
	}

	occupiedSquares := p.getOccupationBitBoard()
	addBlackPawnMoves := func(from coordinate, to coordinate) {
		if to.getRankIndex() != 0 {
			moveCandidates.addMove(from, to, empty)
			return
		}

		moveCandidates.addMove(from, to, blackQueen)
		moveCandidates.addMove(from, to, blackRook)
		moveCandidates.addMove(from, to, blackBishop)
		moveCandidates.addMove(from, to, blackKnight)
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
		moveCandidates.addMove(fromCoord, toCoord, empty)
	}
}

func (p Position) isAttackedByEnemy(player colour, theCoord coordinate) bool {
	return p.controlByColour[player.getOpponent()].get(theCoord)
}

func (p Position) isOccupiedByFriendly(player colour, theCoord coordinate) bool {
	return p.occupationByColour[player].get(theCoord)
}

func (p Position) getOccupationBitBoard() bitBoard {
	return p.occupationByColour[white] | p.occupationByColour[black]
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

	p.activeColour = p.activeColour.getOpponent()

	// TODO check if the player who moved is now in check
	resultsInCheck = false

	return resultsInCheck
}
