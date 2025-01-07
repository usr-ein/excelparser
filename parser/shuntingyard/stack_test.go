package shuntingyard

import (
	"testing"
)

func TestStackPop(t *testing.T) {
	stack := NewStack[int]()
	for i := 0; i < 5; i++ {
		stack.Push(i)
	}

	for i := 4; i > -1; i-- {
		v, ok := stack.Pop()
		if !ok {
			t.Errorf("stack.Pop() failed")
		}
		if v != i {
			t.Errorf("stack.Pop() failed: expected %d, got %d", i, v)
		}
	}

	_, ok := stack.Pop()
	if ok {
		t.Errorf("stack.Pop() succeeded when it should have failed")
	}
}

func TestStackPop2(t *testing.T) {
	stack := NewStack[int]()
	stack.Push(1)
	stack.Push(2)

	if v, ok := stack.Pop(); !ok || v != 2 {
		t.Errorf("stack.Pop() failed: expected 2, got %d", v)
	}

	stack.Push(3)

	if v, ok := stack.Pop(); !ok || v != 3 {
		t.Errorf("stack.Pop() failed: expected 3, got %d", v)
	}

	if v, ok := stack.Pop(); !ok || v != 1 {
		t.Errorf("stack.Pop() failed: expected 1, got %d", v)
	}

	if _, ok := stack.Pop(); ok {
		t.Errorf("stack.Pop() succeeded when it should have failed")
	}
}

func TestStackTop(t *testing.T) {
	stack := NewStack[int]()
	stack.Push(1)
	stack.Push(2)

	v, ok := stack.Top()
	if !ok {
		t.Errorf("stack.Top() failed")
	}
	if v != 2 {
		t.Errorf("stack.Top() failed: expected 2, got %d", v)
	}

	v, ok = stack.Top()
	if !ok {
		t.Errorf("stack.Top() failed")
	}
	if v != 2 {
		t.Errorf("stack.Top() failed: expected 2, got %d", v)
	}

	_, ok = stack.Pop()
	if !ok {
		t.Errorf("stack.Pop() failed")
	}

	v, ok = stack.Top()
	if !ok {
		t.Errorf("stack.Top() failed")
	}
	if v != 1 {
		t.Errorf("stack.Top() failed: expected 1, got %d", v)
	}

	_, ok = stack.Pop()
	if !ok {
		t.Errorf("stack.Pop() failed")
	}

	_, ok = stack.Top()
	if ok {
		t.Errorf("stack.Top() succeeded when it should have failed")
	}
}
