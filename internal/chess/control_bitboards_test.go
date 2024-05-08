package chess

import "testing"

func TestKingControlBitBoards(t *testing.T) {
	result := kingControlFrom[a1]
	expected := bitBoard(0b1100000010)

	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
