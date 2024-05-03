package main

import "testing"

func TestSet(t *testing.T) {
	theBitBoard := bitBoard(0)
	theBitBoard.set(e4)
	expected := bitBoard(1 << e4)

	if theBitBoard != expected {
		t.Errorf("expected %v, got %v", expected, theBitBoard)
	}
}

func TestGet(t *testing.T) {
	theBitBoard := bitBoard(1 << c6)

	if theBitBoard.get(c6) != true {
		t.Errorf("the c6 square should be set")
	}

	if theBitBoard.get(h3) != false {
		t.Errorf("the h3 square should not be set")
	}
}

func TestClear(t *testing.T) {
	theBitBoard := bitBoard(1 << d5)
	theBitBoard.clear(d5)

	if theBitBoard != 0 {
		t.Errorf("the d5 square should be cleared")
	}
}
