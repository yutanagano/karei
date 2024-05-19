package chess

import (
	"fmt"
	"reflect"
)

type algebraicMove struct {
	From      coordinate
	To        coordinate
	Promotion squareState
}

func algebraicMoveFromString(s string) (algebraicMove, error) {
	var result algebraicMove
	var fromSquare, toSquare coordinate
	promotion := empty

	strLen := len(s)

	if !(strLen == 4 || strLen == 5) {
		return result, fmt.Errorf("invalid move: %s", s)
	}

	fromSquare, err := coordinateFromString(s[:2])
	if err != nil {
		return result, err
	}

	toSquare, err = coordinateFromString(s[2:4])
	if err != nil {
		return result, err
	}

	if strLen == 5 {
		if r := toSquare.getRankIndex(); !(r == 0 || r == 7) {
			return result, fmt.Errorf("cannot promote on rank %v: %s", r, s)
		}

		promotion, err = squareStateFromRune(rune(s[4]))
		if err != nil {
			return result, err
		}
	}

	result.From = fromSquare
	result.To = toSquare
	result.Promotion = promotion

	return result, nil
}

func (a algebraicMove) toString() string {
	if a.Promotion != empty {
		return a.From.toString() + a.To.toString() + string(a.Promotion.toRune())
	}

	return a.From.toString() + a.To.toString()
}

func (a algebraicMove) getOffset() int {
	return int(a.To) - int(a.From)
}

type move uint32

const (
	moveOffsetTo             = 6
	moveOffsetCapturedPiece  = moveOffsetTo + 6
	moveOffsetPromotionTo    = moveOffsetCapturedPiece + 4
	moveOffsetCastlingRights = moveOffsetPromotionTo + 4
	moveOffsetEPSquare       = moveOffsetCastlingRights + 4
)

const (
	moveMaskFrom           move = 0x3f
	moveMaskTo             move = 0x3f << moveOffsetTo
	moveMaskCapturedPiece  move = 0xf << moveOffsetCapturedPiece
	moveMaskPromotionTo    move = 0xf << moveOffsetPromotionTo
	moveMaskCastlingRights move = 0xf << moveOffsetCastlingRights
	moveMaskEPSquare       move = 0x7f << moveOffsetEPSquare
)

func moveFromParts(from, to coordinate, capturedPiece, promotionTo squareState, currentCastlingRights castlingRights, currentEPSquare coordinate) move {
	fromEncoded := move(from)
	toEncoded := move(to) << moveOffsetTo
	captureEncoded := move(capturedPiece) << moveOffsetCapturedPiece
	promotionEncoded := move(promotionTo) << moveOffsetPromotionTo
	castlingRightsEncoded := move(currentCastlingRights) << moveOffsetCastlingRights
	EPSquareEncoded := move(currentEPSquare) << moveOffsetEPSquare

	return fromEncoded | toEncoded | captureEncoded | promotionEncoded | castlingRightsEncoded | EPSquareEncoded
}

func (m move) getFromCoordinate() coordinate {
	return coordinate(m & moveMaskFrom)
}

func (m move) getToCoordinate() coordinate {
	return coordinate((m & moveMaskTo) >> moveOffsetTo)
}

func (m move) getCapturedPiece() squareState {
	return squareState((m & moveMaskCapturedPiece) >> moveOffsetCapturedPiece)
}

func (m move) getPromotionTo() squareState {
	return squareState((m & moveMaskPromotionTo) >> moveOffsetPromotionTo)
}

func (m move) getCurrentCastlingRights() castlingRights {
	return castlingRights((m & moveMaskCastlingRights) >> moveOffsetCastlingRights)
}

func (m move) getCurrentEPSquare() coordinate {
	return coordinate((m & moveMaskEPSquare) >> moveOffsetEPSquare)
}

type moveList []algebraicMove

func (l *moveList) addMove(from coordinate, to coordinate, promotion squareState) {
	*l = append(*l, algebraicMove{from, to, promotion})
}

func (l *moveList) filter(evaluator func(algebraicMove) bool) {
	writeIndex := 0

	for _, theMove := range *l {
		if evaluator(theMove) {
			(*l)[writeIndex] = theMove
			writeIndex++
		}
	}

	(*l) = (*l)[:writeIndex]
}

func (l moveList) contains(query algebraicMove) bool {
	for _, element := range l {
		if reflect.DeepEqual(query, element) {
			return true
		}
	}
	return false
}
