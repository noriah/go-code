package slice

import "testing"

func TestSliceStack(t *testing.T) {
	stack := &Stack{}
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	if size := stack.Size(); size != 3 {
		t.Errorf("expected size %d, got %d", 3, size)
	}

	value, err := stack.Peek()
	if err != nil {
		t.Error(err)
	}

	if value.(int) != 3 {
		t.Errorf("expected %d, got %d", 3, value)
	}

	stackPopHelper(t, stack, 3)
	stackPopHelper(t, stack, 2)
	stackPopHelper(t, stack, 1)
}

func stackPopHelper(t *testing.T, stack *Stack, expect int) {
	value, err := stack.Pop()
	if err != nil {
		t.Error(err)
	}

	if expect != value.(int) {
		t.Errorf("expected: %d, got %d", expect, value)
	}
}
