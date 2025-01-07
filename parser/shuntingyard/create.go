package shuntingyard

type ShuntingYardState[T any] struct {
	Operands  Stack[T]
	Operators Stack[Operator]
}

func NewShuntingYardState[T any]() ShuntingYardState[T] {
	operators := NewStack[Operator]()
	operators.Push(Sentinel)
	return ShuntingYardState[T]{
		Operands:  NewStack[T](),
		Operators: operators,
	}
}
