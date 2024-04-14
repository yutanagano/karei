package main

import "testing"

func TestPopFromQueue(t *testing.T) {
	queue := []string{"a", "b", "c"}
	expected_result := "a"
	expected_queue := []string{"b", "c"}
	result := popFromQueue(&queue)
	if result != expected_result {
		t.Errorf("Expected pop result %s, got %s", expected_result, result)
	}
	if len(queue) != len(expected_queue) {
		t.Errorf("Expected return queue to have length %d, got %d", len(expected_queue), len(queue))
	}
}
