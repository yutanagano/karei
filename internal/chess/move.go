package chess

import "fmt"

type Move struct {
	From      coordinate
	To        coordinate
	Promotion squareState
}

func MoveFromString(s string) (Move, error) {
	var result Move
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

func (m Move) ToString() string {
	if m.Promotion != empty {
		return m.From.toString() + m.To.toString() + string(m.Promotion.toRune())
	}

	return m.From.toString() + m.To.toString()
}
