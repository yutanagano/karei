package chess

import "testing"

func TestIsSet(t *testing.T) {
	var cr castlingRights = 0

	if cr.isSet(whiteCastleKingSide) {
		t.Error("whiteCastleKingSide has not been set but isSet returns true")
	}

	cr = whiteCastleQueenSide
	if !cr.isSet(whiteCastleQueenSide) {
		t.Error("whiteCastleQueenSide should be set but isSet returns false")
	}
}

func TestTurnOn(t *testing.T) {
	var cr castlingRights = 0
	cr.turnOn(blackCastleKingSide)

	if cr != blackCastleKingSide {
		t.Error("turning on blackCastleKingSide did not work")
	}
}

func TestTurnOff(t *testing.T) {
	var cr castlingRights = blackCastleQueenSide
	cr.turnOff(blackCastleQueenSide)

	if cr != 0 {
		t.Error("turning off blackCastleQueenSide did not work")
	}
}
