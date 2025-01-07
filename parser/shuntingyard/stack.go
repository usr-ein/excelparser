package shuntingyard

type Stack[T any] interface {
	Push(T)
	Pop() (T, bool)
	Top() (T, bool)
}

type StackImpl[T any] struct {
	data []T
}

func NewStack[T any]() Stack[T] {
	return &StackImpl[T]{}
}

func (s *StackImpl[T]) Push(v T) {
	s.data = append(s.data, v)
}

func (s *StackImpl[T]) Pop() (T, bool) {
	v, ok := s.Top()
	if !ok {
		return v, false
	}
	s.data = s.data[:len(s.data)-1]
	return v, true
}

func (s StackImpl[T]) Top() (T, bool) {
	var v T
	if len(s.data) == 0 {
		return v, false
	}
	return s.data[len(s.data)-1], true
}
