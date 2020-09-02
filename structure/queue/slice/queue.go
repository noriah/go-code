// Package slice implements a Slice Queue.
// A slice queue uses an internal array to hold the values
package slice

import (
	"errors"
	"sync"
)

const defaultSliceSize = 16

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
	// array to hold the data
	array []interface{}
	// Channel to hold our indexes for poping
	popChannel chan int
	// channel to hold our indexes for pushing
	pushChannel chan int
	// Size to keep track how many items are in the channel
	size int
	// the next index to pop
	nextPop int
	// the next index at the end of the array we can push to
	nextPush int
}

// New returns a new Channel Queue
func New() *Queue {
	// Make a new queue
	var newQueue = &Queue{
		array:       make([]interface{}, defaultSliceSize),
		popChannel:  make(chan int, defaultSliceSize),
		pushChannel: make(chan int, defaultSliceSize),
	}

	// Return the queue
	return newQueue
}

func (q *Queue) expand() {
	var newCap = (cap(q.array) + 1) * 2

	var newArray = make([]interface{}, newCap)

	copy(newArray, q.array)

	q.array = newArray

	var newPopChannel = make(chan int, newCap)

	var newPushChannel = make(chan int, newCap)

	var idx int

	// For, or until we somehow break the loop
	for {
		// Select between two options. Pick the first available
		select {
		case idx = <-q.popChannel:
			newPopChannel <- idx

		case idx = <-q.pushChannel:
			newPushChannel <- idx

		default:
			close(q.popChannel)
			q.popChannel = newPopChannel
			close(q.pushChannel)
			q.pushChannel = newPushChannel
			return
		}
	}
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

// Clear removes all items from the queue
func (q *Queue) Clear() {
	// Lock the mutex so we can empty the channel in peace
	q.mu.Lock()

	defer q.mu.Unlock()

	// Set the size to 0
	q.size = 0

	q.nextPop = -1

	q.nextPush = 0

	// For, or until we somehow break the loop
	for {
		// Select between two options. Pick the first available
		select {
		case <-q.popChannel:

		case <-q.pushChannel:

		default:
			return
		}
	}
}

// Push adds a value to the internal channel.
func (q *Queue) Push(value interface{}) {
	var idx int
	var ok bool

	// Lock the mutex so we can push in peace
	q.mu.Lock()

	if q.size >= cap(q.array) {
		q.expand()
	}

	// Select an action
	select {
	case idx, ok = <-q.pushChannel:
		if !ok {
			panic("What happened here with push channel??")
		}

	default:
		idx = q.nextPush
		q.nextPush++
	}

	q.size++

	q.array[idx] = value

	select {
	case q.popChannel <- idx:

	default:
		panic("What??? full pop channel?")
	}

	q.mu.Unlock()
}

// Pop removes a value from the internal channel and returns it.
// Returns nil and error if queue is empty.
func (q *Queue) Pop() (interface{}, error) {
	if q.size == 0 {
		return nil, errorQueueEmpty
	}

	// Lock the mutex so we can pop in peace
	q.mu.Lock()
	// Defer the unlock to after the func exits
	defer q.mu.Unlock()

	var idx = q.nextPop

	q.nextPop = -1

	if idx < 0 {
		var ok bool

		// Select an action
		select {
		// Try to take an item from the channel, and let us know the close state
		case idx, ok = <-q.popChannel:
			// If the channel is not closed
			if !ok {
				panic("What happened here with push channel??")
			}
			q.size--
		default:
			return nil, errorQueueEmpty
		}
	}

	select {
	case q.pushChannel <- idx:

	default:
		panic("What??? full push channel?")
	}

	// How did we get here?
	// Return whatever we got and channel closed error
	return q.array[idx], nil
}

// Peek returns the value at the front of the queue.
// The queue array is not moified
func (q *Queue) Peek() (interface{}, error) {
	if q.size == 0 {
		return nil, errorQueueEmpty
	}

	// Lock the mutex so we can pop in peace
	q.mu.Lock()
	// Defer the unlock to after the func exits
	defer q.mu.Unlock()

	var idx = q.nextPop

	if idx < 0 {
		var ok bool

		// Select an action
		select {
		// Try to take an item from the channel, and let us know the close state
		case idx, ok = <-q.popChannel:
			// If the channel is not closed
			if !ok {
				panic("What happened here with push channel??")
			}
			q.nextPop = idx
		default:
			return nil, errorQueueEmpty
		}
	}

	// How did we get here?
	// Return whatever we got and channel closed error
	return q.array[idx], nil
}
