package util

import (
	"reflect"
	"testing"
)

func TestPopFromQueue(t *testing.T) {
	queue := Queue[string]{"a", "b", "c"}
	expected_result := "a"
	expected_queue := Queue[string]{"b", "c"}
	result := queue.Pop()
	if result != expected_result {
		t.Errorf("expected pop result %s, got %s", expected_result, result)
	}
	if !reflect.DeepEqual(queue, expected_queue) {
		t.Errorf("expected return queue to have length %d, got %d", len(expected_queue), len(queue))
	}
}
