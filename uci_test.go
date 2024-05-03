package main

import (
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

func getMilisecondTimeOutChannel(k uint) chan bool {
	timeOut := make(chan bool, 1)
	go func() {
		time.Sleep(time.Duration(k) * time.Millisecond)
		timeOut <- true
	}()
	return timeOut
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
		{
			"isready",
			"isready",
			[]string{"readyok"},
		},
	}

	fromUci, toUci := startTestUci()

	timeOut := getMilisecondTimeOutChannel(1000)
waitUntilListen:
	for {
		select {
		case output := <-fromUci:
			if output == "info string listening" {
				break waitUntilListen
			}
		case <-timeOut:
			t.Errorf("timeout waiting for UCI to begin listening")
			return
		}
	}

	runTest := func(t *testing.T, test testSpec) {
		toUci <- test.input
		for _, expectedOutput := range test.expectedOutputs {
			timeOut := getMilisecondTimeOutChannel(1000)
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

	type verificationSpec struct {
		squareChecks    []squareCheck
		enPassantSquare coordinate
		castlingRights
		activeColour  colour
		halfMoveClock uint8
	}

	type testSpec struct {
		name  string
		input string
		verificationSpec
	}

	tests := []testSpec{
		{
			"fen",
			"position fen 3rkb1r/p2nqppp/5n2/1B2p1B1/4P3/1Q6/PPP2PPP/2KR3R w k - 3 13",
			verificationSpec{
				[]squareCheck{{d8, blackRook}, {c1, whiteKing}, {g7, blackPawn}, {g5, whiteBishop}, {e1, empty}},
				nullCoordinate,
				0b0100,
				white,
				13,
			},
		},
	}

	_, toUci := startTestUci()

	runTest := func(t *testing.T, test testSpec) {
		toUci <- test.input
		for _, sc := range test.verificationSpec.squareChecks {
			if result := currentPosition.board[sc.coordinate]; result != sc.squareState {
				t.Errorf("expected %v at %v, got %v", sc.squareState, sc.coordinate, result)
			}
		}
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) { runTest(t, test) })
	}
}
