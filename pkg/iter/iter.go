package iter

func New[T any](data []T) Iterator[T] {
	return Iterator[T]{
		data: data,
	}
}

type Iterator[T any] struct {
	data   []T
	offset int
}

func (i *Iterator[T]) Next() (val T, ok bool) {
	if i.offset >= len(i.data) {
		return val, false
	}

	val = i.data[i.offset]
	i.offset++

	return val, true
}
