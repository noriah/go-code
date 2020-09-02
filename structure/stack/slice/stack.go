// Package slice holds implementation for a slice stack.
// A slice stack is a stack implemented with a slice
package slice

import (
	"errors"
	"sync"
)

const defaultSliceSize = 16

// An error to be returned when Pop/Peek-ing on an empty stack
var errorStackEmpty = errors.New("empty stack")

// Stack implements a Slice Stack
type Stack struct {

	// our mutex. don't embed because we don't want to expose it
	mu sync.Mutex

	// the data storage for the stack
	array []interface{}

	// our number of items in the stack. updated on every push and pop
	count int
}

// New returns a new slice Stack
func New(values ...interface{}) *Stack {

	// Make a stack object
	var newStack = &Stack{}

	newStack.array = make([]interface{}, defaultSliceSize)

	// Add any values we may have been passed to the stack
	newStack.Append(values)

	// Return the new stack
	return newStack
}

func (s *Stack) expand() {
	var newArray = make([]interface{}, (cap(s.array)+1)*2)

	copy(newArray, s.array)

	s.array = newArray
}

// Push adds a value to the top of the stack.
//
// Time: O(1) | O(n)
// Space: O(0) | O(n)
func (s *Stack) Push(value interface{}) {

	// Lock the mutex while we are modifying the stack. Prevents someone
	// Adding an item before we do, and having a messed up stack counter
	s.mu.Lock()

	if s.count >= cap(s.array) {
		s.expand()
	}

	s.array[s.count] = value

	// Increment the total items in stack
	s.count++

	// Unlock the mutex
	s.mu.Unlock()
}

// Append adds a list of items to the stack
//
// Time: O(n)
// Space: O(n)
func (s *Stack) Append(values ...interface{}) {

	// assign a variable so we don't do multiple length checks
	var vLen = len(values)

	// If we have less than 2 items in values, we don't want to do the insertion
	// logic below. We only want to check for a single value now and if so, push it.
	if vLen < 2 {

		// Check for length == 1
		if vLen == 1 {

			// If we only have one item in values, just push it.
			s.Push(values[0])
		}

		// end the function
		return
	}

	s.mu.Lock()

	if s.count+vLen >= cap(s.array) {
		s.expand()
	}

	copy(s.array[s.count:s.count+vLen], values)

	s.count += vLen

	// Unlock the mutex
	s.mu.Unlock()
}

// Pop returns the value on the top of the stack, removing it from the stack
//
// Time: O(1)
func (s *Stack) Pop() (interface{}, error) {

	if s.count <= 0 {
		return nil, errorStackEmpty
	}

	s.mu.Lock()

	defer s.mu.Unlock()

	// decrement our count of items in stack
	s.count--

	// return the value in our temp node
	return s.array[s.count+1], nil
}

// Peek returns the value at the front of the stack.
// The stack is not modified.
//
// Time: O(1)
func (s *Stack) Peek() (interface{}, error) {

	if s.count <= 0 {
		return nil, errorStackEmpty
	}

	// Lock the internal mutex to prevent someone pop-ing while we are peek-ing
	s.mu.Lock()

	defer s.mu.Unlock()

	return s.array[s.count], nil
}

// Clear empties the stack.
func (s *Stack) Clear() {

	// Lock our mutex so we can be sure to clear the stack before any other
	// operations happen on it
	s.mu.Lock()

	// Update count to be 0
	s.count = 0

	// Unlock the mutex
	s.mu.Unlock()
}

// Size returns the number of items in the stack
func (s *Stack) Size() int {

	// Lock the mutex so we don't check in the middle of an operation
	s.mu.Lock()

	// To avoid assigning a temp variable just for the count, we can defer
	// the mutex unlock to after we have returned
	defer s.mu.Unlock()

	// return the count of items
	return s.count
}

// IsEmpty checks for stack emptiness
func (s *Stack) IsEmpty() bool {

	// If count is 0, then we have an empty stack
	return s.count == 0
}
