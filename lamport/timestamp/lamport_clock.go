package timestamp

import (
	"fmt"
	"sync"
)

// LamportClock holds the state of a Lamport clock.
type LamportClock struct {
	counter int
}

func (u LamportClock) DefaultTsCmp(v DistributedClock[int]) TsOrder {
	if u.counter > v.Value() {
		return BEF
	} else if u.counter == v.Value() {
		return CON
	} else {
		return AFT
	}
}

func (u LamportClock) Value() int {
	return u.counter
}

func (u LamportClock) String() string {
	return fmt.Sprintf("Lamport{%v}", u.counter)
}

func (clock *LamportClock) Increment() {
	clock.counter++
}

func (clock *LamportClock) Set(data int) {
	clock.counter = data
}

/*
--------------------------
*/

type LamportLocalClock struct {
	clock DistributedClock[int]
	mutex sync.Mutex
}

// Increment the clock by 1 (local event).
func (lc *LamportLocalClock) Forward() {
	lc.mutex.Lock()
	lc.clock.Increment()
	lc.mutex.Unlock()
}

// Adjust compares the current clock value with the received one, and updates the current clock with the maximum of the two, incremented by 1 (message receive event).
func (lc *LamportLocalClock) Adjust(received DistributedClock[int]) error {
	lc.mutex.Lock()
	if received.Value() > lc.clock.Value() {
		lc.clock.Set(received.Value())
	}
	lc.clock.Increment()
	lc.mutex.Unlock()
	return nil
}

// GetTime gets the current time from the Lamport clock.
func (lc *LamportLocalClock) Snapshot() int {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()
	return lc.clock.Value()
}
