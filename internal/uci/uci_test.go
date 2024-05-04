package uci

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/yutanagano/karei/internal/chess"
	"github.com/yutanagano/karei/internal/util"
)

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

	fromUCI, toUCI := startUCIWithDummyEngine()

	err := waitUntilReady(fromUCI, toUCI, 10)
	if err != nil {
		t.Error(err.Error())
	}

	runTest := func(t *testing.T, test testSpec) {
		toUCI <- test.input
		for _, expectedOutput := range test.expectedOutputs {
			timeOut := getMillisecondTimeOutChannel(10)
			select {
			case result := <-fromUCI:
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

func TestHandlePosition(t *testing.T) {
	type testSpec struct {
		name      string
		arguments util.Queue[string]
		fen       chess.FEN
	}

	tests := []testSpec{
		{
			"scotch",
			util.Queue[string]{
				"fen",
				"r1bqkbnr/pppp1ppp/2n5/4p3/3PP3/5N2/PPP2PPP/RNBQKB1R",
				"b",
				"KQkq",
				"-",
				"0",
				"3",
			},
			chess.FEN{
				BoardState:      "r1bqkbnr/pppp1ppp/2n5/4p3/3PP3/5N2/PPP2PPP/RNBQKB1R",
				ActiveColour:    "b",
				CastlingRights:  "KQkq",
				EnPassantSquare: "-",
				HalfMoveClock:   "0",
				FullMoveNumber:  "3",
			},
		},
	}

	runTest := func(t *testing.T, test testSpec) {
		handlePosition(test.arguments)
		expected := chess.Position{}
		expected.LoadFEN(test.fen)

		if !reflect.DeepEqual(currentPosition, expected) {
			t.Errorf("expected position %v, got %v", expected, currentPosition)
		}
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) { runTest(t, test) })
	}
}

func startUCIWithDummyEngine() (fromUCI, toUCI chan string) {
	fromUCI = make(chan string, 100)
	toUCI = make(chan string)
	ConnectClient(toUCI, fromUCI)

	fromDummyEngine := make(chan string)
	toDummyEngine := make(chan string)
	ConnectEngine(fromDummyEngine, toDummyEngine)

	go Start()
	return
}

func waitUntilReady(fromUCI, toUCI chan string, milliseconds uint) error {
	toUCI <- "isready"
	timeOut := getMillisecondTimeOutChannel(milliseconds)
	for {
		select {
		case output := <-fromUCI:
			if output == "readyok" {
				return nil
			}
		case <-timeOut:
			return fmt.Errorf("uci not ready in time (%v ms)", milliseconds)
		}
	}
}

func getMillisecondTimeOutChannel(k uint) chan bool {
	timeOut := make(chan bool, 1)
	go func() {
		time.Sleep(time.Duration(k) * time.Millisecond)
		timeOut <- true
	}()
	return timeOut
}
