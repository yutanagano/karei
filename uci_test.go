package main

import (
	"reflect"
	"testing"
	"time"
)

var fromTestTell []string

func testTell(tokens ...string) {
	output := ""
	for _, token := range tokens {
		output += token
	}
	fromTestTell = append(fromTestTell, output)
}

func TestUci(t *testing.T) {
	tell = testTell
	toUci := make(chan string)
	go uci(toUci)

	for {
		if len(fromTestTell) > 0 && fromTestTell[len(fromTestTell)-1] == "info string listening" {
			break
		}
	}

	tests := []struct {
		name           string
		input          string
		expectedOutput []string
	}{
		{
			"uci",
			"uci",
			[]string{"id name Karei", "id author Yuta Nagano", "option name Hash type spin default 32 min 1 max 1024", "option name Threads type spin default 1 min 1 max 16", "uciok"},
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				fromTestTell = []string{}
				toUci <- test.input
				time.Sleep(10 * time.Millisecond)
				if !reflect.DeepEqual(fromTestTell, test.expectedOutput) {
					t.Errorf("Expected output %v, got %v", test.expectedOutput, fromTestTell)
				}
			},
		)
	}
}

func TestPopFromQueue(t *testing.T) {
	queue := []string{"a", "b", "c"}
	expected_result := "a"
	expected_queue := []string{"b", "c"}
	result := popFromQueue(&queue)
	if result != expected_result {
		t.Errorf("Expected pop result %s, got %s", expected_result, result)
	}
	if !reflect.DeepEqual(queue, expected_queue) {
		t.Errorf("Expected return queue to have length %d, got %d", len(expected_queue), len(queue))
	}
}
