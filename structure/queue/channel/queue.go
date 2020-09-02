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
	mu       sync.Mutex
	channel  chan interface{}
	size     int
	capacity int
}

// New returns a new Channel Queue
// Capacity is the maximum number of items in the queue
func New(capacity int) *Queue {
	var newQueue = &Queue{
		channel:  make(chan interface{}, capacity),
		capacity: capacity,
	}

	return newQueue
}

// Size returns the current number of items in the queue
func (q *Queue) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.size
}

// Capacity returns the maximum number of items in the queue
func (q *Queue) Capacity() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.capacity
}

// Clear removes all items from the queue
func (q *Queue) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.size = 0

	for {
		select {
		case <-q.channel:

		default:
			return
		}
	}
}

// Push adds a value to the internal channel.
// Returns error if queue is full
func (q *Queue) Push(value interface{}) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	select {
	case q.channel <- value:
		q.size++
		return nil
	default:
		return errorQueueFull
	}
}

// Pop removes a value from the internal channel and returns it.
// Returns nil and error if queue is empty.
func (q *Queue) Pop() (interface{}, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	var value interface{}
	var ok bool

	select {
	case value, ok = <-q.channel:
		if ok {
			q.size--
			return value, nil
		}
	default:
		return nil, errorQueueEmpty
	}

	return value, errorChannelClosed
}
