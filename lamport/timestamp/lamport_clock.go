package timestamp

import (
	"fmt"
	"sync"
)

// type DistributedClock[T any] interface {
// 	DefaultTsCmp(DistributedClock[T]) TsOrder
// 	Value() T
// }

// type LocalClock[T any] interface {
// 	Forward()
// 	Adjust(DistributedClock[T]) error
// 	Snapshot() T
// }

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

type LamportLocalClock struct {
	clock LamportClock
	mutex sync.Mutex
}

// Increment the clock by 1 (local event).
func (lc *LamportLocalClock) Forward() {
	lc.mutex.Lock()
	lc.clock.counter++
	lc.mutex.Unlock()
}

// Adjust compares the current clock value with the received one, and updates the current clock with the maximum of the two, incremented by 1 (message receive event).
func (lc *LamportLocalClock) Adjust(received LamportClock) error {
	lc.mutex.Lock()
	if received.Value() > lc.clock.counter {
		lc.clock.counter = received.Value()
	}
	lc.clock.counter++
	lc.mutex.Unlock()
	return nil
}

// GetTime gets the current time from the Lamport clock.
func (lc *LamportLocalClock) Snapshot() int {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()
	return lc.clock.Value()
}
