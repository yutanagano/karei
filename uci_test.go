package main

import (
	"reflect"
	"testing"
	"time"
)

func startTestUci() (chan string, chan string) {
	fromUci := make(chan string)
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

// func TestPosition(t *testing.T) {
// 	type testSpec struct {
// 		squareChecks []struct {
// 			coordinate
// 			squareState
// 		}
// 		enPassantSquare coordinate
// 		castlingRights
// 		activeColour  colour
// 		halfMoveClock uint8
// 	}

// 	tests := []struct {
// 		name string
// 		arguments string
// 		testSpec
// 	}{
// 		{
// 			"fen",
// 			"fen 3rkb1r/p2nqppp/5n2/1B2p1B1/4P3/1Q6/PPP2PPP/2KR3R w k - 3 13",
// 			testSpec{
// 				{ {d8,blackRook}, {c1,whiteKing}, {g7,blackPawn}, {g5,whiteBishop}, {e1,empty} },
// 				nullCoordinate,
// 				0b0100,
// 				white,
// 				13
// 			},
// 		}
// 	}

// 	for _, test := range tests {
// 		t.Run(
// 			test.name,
// 			func (t *testing.T) {
// 				handlePosition(test.arguments)
// 			}
// 		)
// 	}
// }
