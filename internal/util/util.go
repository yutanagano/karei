package util

type Queue[T any] []T

func (q *Queue[T]) Pop() T {
	var result T
	if len(*q) == 0 {
		return result
	}
	result = (*q)[0]
	*q = (*q)[1:]
	return result
}
