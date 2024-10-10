package timestamp

import (
	"sync"
)

// LamportClock holds the state of a Lamport clock.
type LamportClock struct {
	counter int
	mutex   sync.Mutex
}

// Increment the clock by 1 (local event).
func (lc *LamportClock) Increment() {
	lc.mutex.Lock()
	lc.counter++
	lc.mutex.Unlock()
}

// CompareAndUpdate compares the current clock value with the received one, and updates the current clock with the maximum of the two, incremented by 1 (message receive event).
func (lc *LamportClock) CompareAndUpdate(receivedTimestamp int) {
	lc.mutex.Lock()
	if receivedTimestamp > lc.counter {
		lc.counter = receivedTimestamp
	}
	lc.counter++
	lc.mutex.Unlock()
}

// GetTime gets the current time from the Lamport clock.
func (lc *LamportClock) GetTime() int {
	lc.mutex.Lock()
	currentTime := lc.counter
	lc.mutex.Unlock()
	return currentTime
}
