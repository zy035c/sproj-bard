package timestamp

import (
	"encoding/json"
	"fmt"
	"sync"
)

// LamportClock holds the state of a Lamport clock.
type LamportClock struct {
	counter int
}

func (u *LamportClock) DefaultTsCmp(v DistributedClock[int]) TsOrder {
	if u.counter > v.Value() {
		return BEF
	} else if u.counter == v.Value() {
		return CON
	} else {
		return AFT
	}
}

func (u *LamportClock) Value() int {
	return u.counter
}

func (u *LamportClock) String() string {
	return fmt.Sprintf("Lamport{%v}", u.counter)
}

func (clock *LamportClock) Increment() {
	clock.counter++
}

func (clock *LamportClock) Set(data int) {
	clock.counter = data
}

func (clock *LamportClock) Clone() DistributedClock[int] {
	return &LamportClock{
		counter: clock.Value(),
	}
}

func (clock *LamportClock) MarshalJSON() ([]byte, error) {
	return json.Marshal(DistClockJson[int]{Data: clock.counter, ClockType: "lamport"})
}

func (clock *LamportClock) UnmarshalJSON(data []byte) error {
	var dict DistClockJson[int]
	err := json.Unmarshal(data, &dict)
	if err != nil {
		return err
	}
	clock.counter = dict.Data
	return nil
}

var _ DistributedClock[int] = &LamportClock{}

/*
--------------------------
*/

type LamportLocalClock struct {
	Clock DistributedClock[int]
	mutex sync.Mutex
}

// Increment the clock by 1 (local event).
func (lc *LamportLocalClock) Forward() {
	lc.mutex.Lock()
	lc.Clock.Increment()
	lc.mutex.Unlock()
}

// Adjust compares the current clock value with the received one, and updates the current clock with the maximum of the two, incremented by 1 (message receive event).
func (lc *LamportLocalClock) Adjust(received DistributedClock[int]) error {
	lc.mutex.Lock()
	if received.Value() > lc.Clock.Value() {
		lc.Clock.Set(received.Value())
	}
	lc.Clock.Increment()
	lc.mutex.Unlock()
	return nil
}

// GetTime gets the current time from the Lamport clock.
func (lc *LamportLocalClock) Snapshot() int {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()
	return lc.Clock.Value()
}

func (lc *LamportLocalClock) SnapshotTS() DistributedClock[int] {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()
	return lc.Clock.Clone()
}
