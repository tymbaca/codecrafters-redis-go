package option

func Wrap[T any](val T, set bool) Value[T] {
	if set {
		return Some(val)
	}

	return None[T]()
}

func Some[T any](val T) Value[T] {
	return Value[T]{
		val: val,
		set: true,
	}
}

func None[T any]() Value[T] {
	var zero T
	return Value[T]{
		val: zero,
		set: false,
	}
}

type Value[T any] struct {
	val T
	set bool
}

func (v Value[T]) Get() (T, bool) {
	return v.val, v.set
}

func (v Value[T]) IsSet() bool {
	return v.set
}
