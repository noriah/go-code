// Package slice implements a Slice Queue.
// A slice queue uses an internal array to hold the values
package slice

import (
	"errors"
	"sync"
)

const defaultSliceSize = 16

// An error to be returned when Pop-ing on an empty queue
var errorQueueEmpty = errors.New("empty queue")

// An error to be returned if we hit a closed channel (we never should)
var errorChannelClosed = errors.New("closed channel")

// Queue is a Size that is also a queue
// It uses channels internally to keep track of open slots
type Queue struct {
	// Mutex to lock when we are modifying things
	mu sync.Mutex
	// array to hold the data
	array []interface{}
	// Channel to hold our indexes for poping
	popChannel chan int
	// channel to hold our indexes for pushing
	pushChannel chan int
	// Size to keep track how many items are in the array
	size int
	// the next index to pop
	nextPop int
	// the next index at the end of the array we can push to
	nextPush int
}

// New returns a new Slice Queue
func New() *Queue {
	// Make a new queue
	var newQueue = &Queue{
		array:       make([]interface{}, defaultSliceSize),
		popChannel:  make(chan int, defaultSliceSize),
		pushChannel: make(chan int, defaultSliceSize),
		nextPop:     -1,
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
	// Lock the mutex so we can empty the channels in peace
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

// IsEmpty returns the emptiness state
func (q *Queue) IsEmpty() bool {
	return q.Size() == 0
}

// Push adds a value to the internal array.
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
		default:
			return nil, errorQueueEmpty
		}
	}

	select {
	case q.pushChannel <- idx:
		q.size--

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
