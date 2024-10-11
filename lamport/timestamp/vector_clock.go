package timestamp

import (
	"fmt"
	"lamport/utils"
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

type VectorClock struct {
	vector []uint64
	idx    uint32
	size   uint32
}

func NewVectorClock(sz uint32, idx uint32) (*VectorClock, error) {
	if idx >= sz {
		return nil, fmt.Errorf("idx %d is eq or gt than size %d", idx, sz)
	}
	vector := make([]uint64, sz)
	return &VectorClock{
		vector: vector,
		idx:    idx,
		size:   sz,
	}, nil
}

func (clock VectorClock) DefaultTsCmp(other DistributedClock[[]uint64]) TsOrder {
	n_bef := 0
	n_aft := 0

	v := clock.Value()
	u := other.Value()

	if len(v) != len(u) {
		return CON
	}

	for i := 0; i < len(v); i++ {
		if u[i] > v[i] {
			n_aft++
		}
		if u[i] < v[i] {
			n_bef++
		}
	}
	fmt.Printf("n_aft %v n_bef %v\n", n_aft, n_bef)
	if n_bef != 0 && n_aft != 0 {
		return CON
	} else if n_bef == 0 && n_aft == 0 {
		return CON
	} else if n_bef != 0 {
		return BEF
	} else {
		return AFT
	}
}

func (clock VectorClock) Value() []uint64 {
	return utils.SliceCpy[uint64](clock.vector)
}

type VectorLocalClock struct {
	clock VectorClock
	mutex sync.Mutex
	idx   uint32
	size  uint32
}

func NewVectorLocalClock(sz uint32, idx uint32) (*VectorLocalClock, error) {
	res, err := NewVectorClock(sz, idx)
	if err != nil {
		return nil, err
	}
	return &VectorLocalClock{
		clock: *res,
		size:  sz,
		idx:   idx,
	}, nil
}

func (vtc *VectorLocalClock) Forward() {
	vtc.mutex.Lock()
	defer vtc.mutex.Unlock()
	vtc.clock.vector[vtc.idx]++
}

func (vtc *VectorLocalClock) Adjust(m VectorClock) error {
	vtc.mutex.Lock()
	defer vtc.mutex.Unlock()
	if m.size != vtc.size {
		return fmt.Errorf("ts has a size of %d, local ts has size %d", m.size, vtc.size)
	}

	for i, ts := range vtc.clock.vector {
		if m.vector[i] > ts {
			vtc.clock.vector[i] = m.vector[i]
		}
	}
	vtc.clock.vector[vtc.idx]++
	return nil
}

func (vtc *VectorLocalClock) Snapshot() []uint64 {
	return vtc.clock.Value()
}
