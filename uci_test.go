package main

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func startTestUci() (chan string, chan string) {
	fromUci := make(chan string, 100)
	toUci := make(chan string)

	tell = func(tokens ...string) {
		output := ""
		for _, token := range tokens {
			output += token
		}
		fromUci <- output
	}

	go uci(toUci)

	return fromUci, toUci
}

func getMillisecondTimeOutChannel(k uint) chan bool {
	timeOut := make(chan bool, 1)
	go func() {
		time.Sleep(time.Duration(k) * time.Millisecond)
		timeOut <- true
	}()
	return timeOut
}

func waitUntilReady(fromUci chan string, toUci chan string, milliseconds uint) error {
	toUci <- "isready"
	timeOut := getMillisecondTimeOutChannel(milliseconds)
	for {
		select {
		case output := <-fromUci:
			if output == "readyok" {
				return nil
			}
		case <-timeOut:
			return fmt.Errorf("uci not ready in time (%v ms)", milliseconds)
		}
	}
}

func TestUci(t *testing.T) {
	type testSpec struct {
		name            string
		input           string
		expectedOutputs []string
	}

	tests := []testSpec{
		{
			"uci",
			"uci",
			[]string{
				"id name Karei",
				"id author Yuta Nagano",
				"option name Hash type spin default 32 min 1 max 1024",
				"option name Threads type spin default 1 min 1 max 16",
				"uciok",
			},
		},
		{
			"debug on",
			"debug on",
			[]string{"info string debug mode on"},
		},
		{
			"debug off",
			"debug off",
			[]string{"info string debug mode off"},
		},
	}

	fromUci, toUci := startTestUci()
	err := waitUntilReady(fromUci, toUci, 10)
	if err != nil {
		t.Error(err.Error())
	}

	runTest := func(t *testing.T, test testSpec) {
		toUci <- test.input
		for _, expectedOutput := range test.expectedOutputs {
			timeOut := getMillisecondTimeOutChannel(10)
			select {
			case result := <-fromUci:
				if result != expectedOutput {
					t.Errorf("expected output %v, got %v", expectedOutput, result)
				}
			case <-timeOut:
				t.Errorf("timeout waiting for output %v", expectedOutput)
			}
		}
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) { runTest(t, test) })
	}
}

func TestPopFromQueue(t *testing.T) {
	queue := []string{"a", "b", "c"}
	expected_result := "a"
	expected_queue := []string{"b", "c"}
	result := popFromQueue(&queue)
	if result != expected_result {
		t.Errorf("expected pop result %s, got %s", expected_result, result)
	}
	if !reflect.DeepEqual(queue, expected_queue) {
		t.Errorf("expected return queue to have length %d, got %d", len(expected_queue), len(queue))
	}
}

func TestPosition(t *testing.T) {
	type squareCheck struct {
		coordinate
		squareState
	}

	type testSpec struct {
		name            string
		arguments       []string
		squareChecks    []squareCheck
		enPassantSquare coordinate
		castlingRights
		activeColour  colour
		halfMoveClock uint8
	}

	tests := []testSpec{
		{
			"fen",
			[]string{"fen", "3rkb1r/p2nqppp/5n2/1B2p1B1/4P3/1Q6/PPP2PPP/2KR3R", "w", "k", "-", "3", "13"},
			[]squareCheck{
				{d8, blackRook},
				{c1, whiteKing},
				{g7, blackPawn},
				{g5, whiteBishop},
				{e1, empty},
			},
			nullCoordinate,
			0b0100,
			white,
			3,
		},
	}

	runTest := func(t *testing.T, test testSpec) {
		handlePosition(test.arguments)

		for _, sc := range test.squareChecks {
			if result := currentPosition.board[sc.coordinate]; result != sc.squareState {
				t.Errorf("expected %v at %v, got %v", sc.squareState, sc.coordinate, result)
			}

			if sc.squareState != empty {
				theColour := sc.squareState.getColour()
				thePieceType := sc.squareState.getPieceType()
				if !currentPosition.colourMasks[theColour].get(sc.coordinate) {
					t.Errorf("colourMask not set for %v at %v", sc.squareState, sc.coordinate)
				}
				if !currentPosition.pieceTypeMasks[thePieceType].get(sc.coordinate) {
					t.Errorf("pieceTypeMask not set for %v at %v", sc.squareState, sc.coordinate)
				}
			}
		}

		if currentPosition.enPassantSquare != test.enPassantSquare {
			t.Errorf("expected en passant square %v, got %v", test.enPassantSquare, currentPosition.enPassantSquare)
		}

		if currentPosition.castlingRights != test.castlingRights {
			t.Errorf("expected castling rights %v, got %v", test.castlingRights, currentPosition.castlingRights)
		}

		if currentPosition.activeColour != test.activeColour {
			t.Errorf("should be %v to move", test.activeColour)
		}

		if currentPosition.halfMoveClock != test.halfMoveClock {
			t.Errorf("expected half move clock to be %v, got %v", test.halfMoveClock, currentPosition.halfMoveClock)
		}
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) { runTest(t, test) })
	}
}
