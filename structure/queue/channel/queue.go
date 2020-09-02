// Package channel implements a Channel Queue.
// A channel queue uses an internal channel to hold the values,
// It is fixed in size and cannot be peeked.
package channel

import (
	"errors"
	"sync"
)

// An error to be returned when Push-ing on a full queue
var errorQueueFull = errors.New("full queue")

// An error to be returned when Pop-ing on an empty queue
var errorQueueEmpty = errors.New("empty queue")

// An error to be returned if we hit a closed channel (we never should)
var errorChannelClosed = errors.New("closed channel")

// Queue is a channel that is also a queue but has no peek
// Size is fixed. Adding to a full channel queue will return error
type Queue struct {
	// Mutex to lock when we are modifying things
	mu sync.Mutex
	// Channel to hold our values
	channel chan interface{}
	// Size to keep track how many items are in the channel
	size int
	// Capacity of the channel. Never changes
	capacity int
}

// New returns a new Channel Queue
// Capacity is the maximum number of items in the queue
func New(capacity int) *Queue {
	// Make a new queue
	var newQueue = &Queue{
		// make a channel that accepts interface{}, with size of capacity
		channel: make(chan interface{}, capacity),
		// Set our capacity to capacity
		capacity: capacity,
	}

	// Return the queue
	return newQueue
}

// Size returns the current number of items in the queue
func (q *Queue) Size() int {
	// Lock the mutex so we can get the size at time of the queue
	q.mu.Lock()
	// Defer the unlock to after we return
	defer q.mu.Unlock()

	// Return the size
	return q.size
}

// Capacity returns the maximum number of items in the queue
func (q *Queue) Capacity() int {
	// Return the capacity
	return q.capacity
}

// Clear removes all items from the queue
func (q *Queue) Clear() {
	// Lock the mutex so we can empty the channel in peace
	q.mu.Lock()
	// Defer the unlock to after the func exits
	defer q.mu.Unlock()

	// Set the size to 0
	q.size = 0

	// For, or until we somehow break the loop
	for {
		// Select between two options. Pick the first available
		select {
		// Remove an item from the channel
		case <-q.channel:

			// Default action if we can't take things out of the channel (cuz its empty)
		default:
			// Exit the function
			return
		}
	}
}

// Push adds a value to the internal channel.
// Returns error if queue is full
func (q *Queue) Push(value interface{}) error {
	// Lock the mutex so we can push in peace
	q.mu.Lock()
	// Defer the unlock to after the func exits
	defer q.mu.Unlock()

	// Select an action
	select {
	// Try to send a value on the channel
	case q.channel <- value:
		// We sent! increment the size
		q.size++
		// return no error
		return nil
		// Unable to send on channel
	default:
		// Return queue full error
		return errorQueueFull
	}
}

// Pop removes a value from the internal channel and returns it.
// Returns nil and error if queue is empty.
func (q *Queue) Pop() (interface{}, error) {
	// Lock the mutex so we can pop in peace
	q.mu.Lock()
	// Defer the unlock to after the func exits
	defer q.mu.Unlock()

	// Define some variables for later.
	var value interface{}

	// If we want to have the value, and ok for check, we have to predefine this
	var ok bool

	// Select an action
	select {
	// Try to take an item from the channel, and let us know the close state
	case value, ok = <-q.channel:
		// If the channel is not closed
		if ok {
			// Decrement the size
			q.size--
			// return the value and no error
			return value, nil
		}
		// If we can't pull from the channel
	default:
		// Return nil and Queue empty error
		return nil, errorQueueEmpty
	}

	// How did we get here?
	// Return whatever we got and channel closed error
	return value, errorChannelClosed
}
