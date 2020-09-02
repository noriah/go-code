// Package linked holds implementation for a linked stack.
// A linked stack is a stack implemented with one-way linked nodes
package linked

import (
	"errors"
	"sync"
)

// An error to be returned when Pop/Peek-ing on an empty stack
var errorQueueEmpty = errors.New("empty stack")

// node holds an entry in the stack
type node struct {

	// reference to the next item in a stack
	next *node

	// value held by this node
	value interface{}
}

// Stack implements a Linked Stack
type Stack struct {

	// our mutex. don't embed because we don't want to expose it
	mu sync.Mutex

	// the head node of the stack. updated on every push and pop
	head *node

	// our number of items in the stack
	count int
}

// New returns a new Linked Stack
func New(values ...interface{}) *Stack {

	// We always add items to the top of the stack, so we only care
	// about keeping track of whats at the top. Each node points to the
	// node below it in the stack.
	// This allows us to have simplified logic for the Stack.Push and
	// Stack.Append methods

	// Make a stack object
	newQueue := &Stack{}

	// Add any values we may have been passed to the stack
	newQueue.Append(values)

	// Return the new stack
	return newQueue
}

// Push adds a value to the top of the stack.
//
// Time: O(1)
// Space: O(1)
func (s *Stack) Push(value interface{}) {

	// Lock the mutex while we are modifying the stack. Prevents someone
	// Adding a node before we do, and having a messed up stack
	s.mu.Lock()

	// Make a new node to be added to the stack
	s.head = &node{

		// Set next on our new node to be the current top of our stack.
		next: s.head,

		// Set the value of our new node.
		value: value,
	}

	// Increment the total items in stack
	s.count++

	// Unlock the mutex
	s.mu.Unlock()
}

// Append adds values to the top of the stack by building a mini-stack and
// then updating the stack head to be our mini-stack head.
// The internal mutex is locked once we have built up a collection
// of nodes to append.
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

	// Define variables to build a mini-stack
	var head, last *node

	// Start with the last node and build up
	head = &node{
		// Set the value of our tail node
		value: values[0],
	}

	// save a reference to the last item in our mini-stack so we can point
	// its next value to the top of our current stack
	last = head

	// For all the values left in the array, iterate backwards, building our
	// mini-stack from the bottom up
	for idx := 1; idx < vLen; idx++ {

		// Make a new node and assign it to our variable
		// NOTE: even though we are assigning next to be a new value, the body of
		// the node instantiation is evaluated first, so we don't have to worry about
		// pointing a new node to itself
		head = &node{

			// Set the next value on our new node to be the previous node that we made
			next: head,

			// Set the value of the new node
			value: values[idx],
		}
	}

	// Lock the mutex so we can be sure to add our new mini-stack to the
	// end of the actual stack. Without this, we could be in a race condition
	// where we took long enough to build the mini-stack that another
	s.mu.Lock()

	// Set next on the last mini-stack item to point to our real stack top
	last.next = s.head

	// Set head on our stack to point to the top of our mini-stack
	s.head = head

	// increase our count by number of values
	s.count += vLen

	// Unlock the mutex
	s.mu.Unlock()
}

// Pop returns the value on the top of the stack, removing it from the stack
//
// Time: O(1)
func (s *Stack) Pop() (interface{}, error) {

	// If our tail node is the same our our head node, then we have an empty stack
	if s.head == nil {

		// Return a nil value and our error
		return nil, errorQueueEmpty
	}

	// Define a node pointer to hold the head
	var temp *node

	// Lock the mutex so nobody can modify the stack while we are removing
	// the head of the stack
	s.mu.Lock()

	// assign the current head node to our variable so we don't lose it
	temp = s.head

	// set the stack to point to the next item still in stack
	s.head = temp.next

	// decrement our count of items in stack
	s.count--

	// Unlock the mutex
	s.mu.Unlock()

	// return the value in our temp node
	return temp.value, nil
}

// Peek returns the value at the front of the stack.
// The stack is not modified.
//
// Time: O(1)
func (s *Stack) Peek() (interface{}, error) {

	// Empty stack check
	if s.head == nil {
		return nil, errorQueueEmpty
	}

	// Make a temporary pointer
	var temp *node

	// Lock the internal mutex to prevent someone pop-ing while we are peek-ing
	s.mu.Lock()

	// Set our temp pointer to be the head item in our stack
	temp = s.head

	// Unlock the mutex
	s.mu.Unlock()

	// Return the the value in our temp node
	return temp.value, nil
}

// Clear empties the stack.
// Since the garbage collector cleans up all pointer values once they are no
// longer referenced, we just need to set our tail pointer to our head node,
// and set next on the head node to our tail pointer value (which is our head node).
func (s *Stack) Clear() {

	// Lock our mutex so we can be sure to clear the stack before any other
	// operations happen on it
	s.mu.Lock()

	// Update our head node to point next to our tail
	s.head = nil

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

	// If our head is nil, then we have an empty stack
	return s.head == nil
}
