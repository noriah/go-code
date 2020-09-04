// Package linked holds implementation for a Linked List Queue.
// A linked list queue is implemented with one-way linked nodes
// Each node points to the next node in the queue. The last node points
// To the root (sentinel) node.
package linked

import (
	"errors"
	"sync"
)

// An error to be returned when Push-ing on a full queue
var errorQueueFull = errors.New("full queue")

// An error to be returned when Dequeue/Peek-ing on an empty queue
var errorQueueEmpty = errors.New("empty queue")

// Node is a thin wrapper around a value in a queue.
// It holds the value and a reference to the next node in the queue.
type Node struct {
	next  *Node       // Reference to next node in our queue
	value interface{} // Value this node represents in our queue
}

// Queue implements a Linked List Queue
// References to the head and tail of the queue are held so that we can
// achieve O(1) time for insertion and removal
// The queue is empty when the tail points to our root (root is the same as tail)
type Queue struct {
	mu       sync.Mutex // Mutex for safe parallel operations
	root     *Node      // Root node of our queue. Sentinel node
	tail     *Node      // Tail node of our queue. Real node unless empty, then root
	count    int        // Total number of nodes minus root node
	capacity int        // Maximum size of our queue. 0 means no limit (dynamic)
}

// New returns a new Linked List Queue.
// The optional size may be specified. Only the first value will be used.
func New(size ...int) *Queue {
	var capacity = 0

	if len(size) > 0 {
		capacity = size[0]
		if capacity < 0 {
			panic("Negative value for size provided")
		}
	}

	// Make a queue object
	var newQueue = &Queue{

		// Assign an empty root node
		root: &Node{},

		// set the capacity (if there is any)
		capacity: capacity,
	}

	// Setup our queue
	newQueue.clear()

	// Return the new queue
	return newQueue
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

// IsFull returns the fullness state
func (q *Queue) IsFull() bool {
	// If our queue has a capacity set, honor it
	return q.capacity > 0 && q.count == q.capacity
}

// Clear empties the queue.
// Since the garbage collector cleans up all pointer values once they are no
// longer referenced, we just need to set our tail pointer to our root node,
// and set next on the root node to our tail pointer value (which is our root node).
func (q *Queue) Clear() {

	// Lock our mutex so we can be sure to clear the queue before any other
	// operations happen on it
	q.mu.Lock()

	// Do our clear things
	q.clear()

	// Unlock the mutex
	q.mu.Unlock()
}

// Enqueue inserts a value at the end of the queue.
//
// First check for fullness. If full, return error.
// Make a new node. Set node.next to point to the root of the queue.
// Add value to node.
// Set the the queue tail node next value to point to our new node.
// Set the tail value of our queue to be our new node
// Increment the count of nodes
// Unlock the mutex so other operations can continue
// Return no error
//
// Time: O(1)
// Space: O(1)
func (q *Queue) Enqueue(value interface{}) error {
	// Fullness check
	if q.IsFull() {
		// Return error on full
		return errorQueueFull
	}

	// Make a new node to be added to the queue
	var newNode = &Node{

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

	return nil
}

// Append adds values to the end of the queue by building a mini-list and
// then updating the end of the queue to point to our list.
// The internal mutex is locked once we have built up a collection
// of nodes to append.
//
// Time: O(n)
// Space: O(n)
func (q *Queue) Append(values ...interface{}) {

	// assign a variable so we don't do multiple length checks
	var vLen = len(values)

	// If we have less than 2 items in values, we don't want to do the insertion
	// logic below. We only want to check for a single value now and if so, push it.
	if vLen < 2 {

		// Check for length == 1
		if vLen == 1 {

			// If we only have one item in values, just push it.
			q.Enqueue(values[0])
		}

		// end the function
		return
	}

	// Define variables to build a mini-queue
	var next, tail *Node
	var idx = vLen - 1

	// Make a tail node and build up
	tail = &Node{

		// Set the next on our tail to be the root of the queue.
		next: q.root,

		// Set the value of our tail node
		value: values[idx],
	}

	// Set the next to be the tail of our mini-queu
	next = tail

	// For all the values left in the array, iterate backwards, building our
	// mini-queue from the bottom up
	for idx--; idx >= 0; idx-- {

		// Make a new node and assign it to our variable
		// NOTE: even though we are assigning next to be a new value, the body of
		// the node instantiation is evaluated first, so we don't have to worry about
		// pointing a new node to itself
		next = &Node{

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

// Dequeue returns the value at the front of the queue, removing it from the queue
//
// Time: O(1)
func (q *Queue) Dequeue() (interface{}, error) {

	// If our tail node is the same our our root node, then we have an empty queue
	if q.tail == q.root {

		// Return a nil value and our error
		return nil, errorQueueEmpty
	}

	// Define a node pointer to hold the head
	var temp *Node

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

	// Lock the internal mutex to prevent someone pop-ing while we are peek-ing
	q.mu.Lock()

	// Unlock the mutex
	defer q.mu.Unlock()

	// Return the the value in our temp node
	return q.root.next.value, nil
}

// Helper Methods
// These methods are used internally.

// clear updates the tail and root nodes to be the same, and points the root node
// next to be itself. This is so we can easily
func (q *Queue) clear() {

	// Point the tail of our queue to the root node.
	// The queue is needs to be empty, so we will want to add items to
	// the root node
	// This allows us to have simplified logic
	// for the Queue.Enqueue and Queue.Append methods
	q.tail = q.root

	// Assign the next value on the root node to point to our tail
	q.root.next = q.tail

	// Update count to be 0
	q.count = 0
}
