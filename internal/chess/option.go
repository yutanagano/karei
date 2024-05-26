package chess

type option[T any] struct {
	value      T
	filledFlag bool
}

func newFilledOption[T any](value T) option[T] {
	return option[T]{value: value, filledFlag: true}
}

func newEmptyOption[T any]() option[T] {
	return option[T]{filledFlag: false}
}

func (o option[_]) isEmpty() bool {
	return !o.filledFlag
}

func (o option[T]) getValue() T {
	if o.isEmpty() {
		panic("cannot get value of empty option.")
	}
	return o.value
}

func (o *option[T]) clear() {
	o.filledFlag = false
}

func (o *option[T]) setValue(value T) {
	o.value = value
	o.filledFlag = true
}
