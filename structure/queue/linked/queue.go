// Package linked holds implementation for a linked queue.
// A linked queue is a queue implemented with one-way linked nodes
package linked

import (
	"errors"
	"sync"
)

// An error to be returned when Pop/Peek-ing on an empty queue
var errorQueueEmpty = errors.New("empty queue")

// node holds an entry in the queue
type node struct {
	// reference to the next item in the queue
	next *node

	// value held by this node
	value interface{}
}

// Queue implements a Linked Queue
// References to the head and tail of the queue are held so that we can
// achieve O(1) time for insertion and removal
// The queue is empty when the tail points to our root (root is the same as tail)
type Queue struct {

	// our mutex. don't embed because we don't want to expose it
	mu sync.Mutex

	// the root node of the queue. its next value changes on every pop
	root *node

	// the tail of our queue. updated every push
	tail *node

	// our number of items in the queue
	count int
}

// New returns a new Linked Queue
func New(values ...interface{}) *Queue {

	// Make a queue object
	newQueue := &Queue{
		root: &node{},
	}

	// Point the tail of our queue to the queue itself.
	// The queue is empty right now, so we will want to add items to
	// the current root node.
	// This allows us to have simplified logic
	// for the Queue.Push and Queue.Append methods
	newQueue.tail = newQueue.root

	// Assign the next value on the root node to point to our tail
	newQueue.root.next = newQueue.tail

	// Add any values we may have been passed to the queue
	newQueue.Append(values)

	// Return the new queue
	return newQueue
}

// Push appends a value to the end of the queue.
//
// Time: O(1)
// Space: O(1)
func (q *Queue) Push(value interface{}) {

	// Make a new node to be added to the queue
	newNode := &node{

		// Set next on our new node to be the head of our queue. This allows
		// us to easily add to the queue when it is empty once again.
		next: q.root,

		// Set the value of our new node.
		value: value,
	}

	// Lock the mutex while we are modifying the queue. Prevents someone
	// Adding a node before we do, and having a messed up queue
	q.mu.Lock()

	// Set the next node value at the tail of our queue to be our new node
	q.tail.next = newNode

	// Set the tail of our queue to be our new node
	q.tail = newNode

	// Increment the total items in queue
	q.count++

	// Unlock the mutex
	q.mu.Unlock()
}

// Append adds values to the end of the queue by building a mini-list and
// then updating the end of the queue to point to our list.
// The internal mutex is locked once we have built up a collection
// of nodes to append.
//
// Time: O(n)
// Space: O(n)
func (q *Queue) Append(values ...interface{}) {

	// Only do things if values is not empty, and more than one exists
	if vLen := len(values); vLen > 1 {

		// Define variables to build a mini-queue
		var next, tail *node

		// Make a tail node and build up
		tail = &node{

			// Set the next on our tail to be the root of the queue.
			next: q.root,

			// Set the value of our tail node
			value: values[vLen-1],
		}

		// Set the next to be the tail of our mini-queu
		next = tail

		// For all the values left in the array, iterate backwards, building our
		// mini-queue from the bottom up
		for idx := vLen - 2; idx >= 0; idx-- {

			// Make a new node and assign it to our variable
			// NOTE: even though we are assigning next to be a new value, the body of
			// the node instantiation is evaluated first, so we don't have to worry about
			// pointing a new node to itself
			next = &node{

				// Set the next value on our new node to be the previous node that we made
				next: next,

				// Set the value of the new node
				value: values[idx],
			}
		}

		// Lock the mutex so we can be sure to add our new mini-queue to the
		// end of the actual queue. Without this, we could be in a race condition
		// where we took long enough to build the mini-queue that another
		q.mu.Lock()

		// Set next on the tail queue item to point to our mini-queue start
		q.tail.next = next

		// Set tail on our queue to point to the tail of our mini-queue
		q.tail = tail

		// increase our count by number of values
		q.count += vLen

		// Unlock the mutex
		q.mu.Unlock()
	}
}

// Pop returns the value at the front of the queue, removing it from the queue
//
// Time: O(1)
func (q *Queue) Pop() (interface{}, error) {

	// If our tail node is the same our our root node, then we have an empty queue
	if q.tail == q.root {

		// Return a nil value and our error
		return nil, errorQueueEmpty
	}

	// Define a node pointer to hold the head
	var temp *node

	// Lock the mutex so nobody can modify the queue while we are removing
	// the head of the queue
	q.mu.Lock()

	// assign the current head node to our variable so we don't lose it
	temp = q.root.next

	// set the queue to point to the next item still in queue
	q.root.next = temp.next

	// decrement our count of items in queue
	q.count--

	// Unlock the mutex
	q.mu.Unlock()

	// return the value in our temp node
	return temp.value, nil
}

// Peek returns the value at the front of the queue.
// The queue is not modified.
//
// Time: O(1)
func (q *Queue) Peek() (interface{}, error) {
	// Empty queue check
	if q.tail == q.root {
		return nil, errorQueueEmpty
	}

	// Make a temporary pointer
	var temp *node

	// Lock the internal mutex to prevent someone pop-ing while we are peek-ing
	q.mu.Lock()

	// Set our temp pointer to be the head item in our queue
	temp = q.root.next

	// Unlock the mutex
	q.mu.Unlock()

	// Return the the value in our temp node
	return temp.value, nil
}

// Clear empties the queue.
// Since the garbage collector cleans up all pointer values once they are no
// longer referenced, we just need to set our tail pointer to our root node,
// and set next on the root node to our tail pointer value (which is our root node).
func (q *Queue) Clear() {
	// Lock our mutex so we can be sure to clear the queue before any other
	// operations happen on it
	q.mu.Lock()

	// Update our tail to point to our root node
	q.tail = q.root

	// Update our root node to point next to our tail
	q.root.next = q.tail

	// Unlock the mutex
	q.mu.Unlock()
}

// Size returns the number of items in the queue
func (q *Queue) Size() int {
	// Lock the mutex so we don't check in the middle of an operation
	q.mu.Lock()

	// To avoid assigning a temp variable just for the count, we can defer
	// the mutex unlock to after we have returned
	defer q.mu.Unlock()

	// return the count of items
	return q.count
}

// IsEmpty checks for queue emptiness
func (q *Queue) IsEmpty() bool {
	// If our tail points to our root, then we have an empty queue
	return q.tail == q.root
}
