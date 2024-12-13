package timestamp

import (
	"encoding/json"
	"fmt"
	"lamport/utils"
	"sync"
)

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
	// fmt.Printf("n_aft %v n_bef %v\n", n_aft, n_bef)
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

func (clock *VectorClock) Increment() {
	clock.vector[clock.idx]++
}

func (clock *VectorClock) Set(data []uint64) {
	clock.vector = data
}

func (clock *VectorClock) Clone() DistributedClock[[]uint64] {
	return &VectorClock{
		vector: clock.Value(),
		size:   clock.size,
		idx:    clock.idx,
	}
}

func (clock *VectorClock) String() string {
	return ""
}

func (clock *VectorClock) MarshalJSON() ([]byte, error) {
	return json.Marshal(DistClockJson[[]uint64]{})
}

func (clock *VectorClock) UnmarshalJSON(json []byte) error {
	panic("not implemented! (clock *VectorClock) UnmarshalJSON(json []byte) error")
	// return nil
}

var _ DistributedClock[[]uint64] = &VectorClock{}

/*
--------------------------
*/

type VectorLocalClock struct {
	Clock DistributedClock[[]uint64]
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
		Clock: res,
		size:  sz,
		idx:   idx,
	}, nil
}

func (vtc *VectorLocalClock) Forward() {
	vtc.mutex.Lock()
	defer vtc.mutex.Unlock()
	vtc.Clock.Increment()
}

func (vtc *VectorLocalClock) Adjust(m DistributedClock[[]uint64]) error {
	vtc.mutex.Lock()
	defer vtc.mutex.Unlock()
	mlen := len(m.Value())
	if mlen != int(vtc.size) {
		return fmt.Errorf("ts has a size of %d, local ts has size %d", mlen, vtc.size)
	}

	vec := vtc.Clock.Value()
	for i, ts := range vec {
		if m.Value()[i] > ts {
			vec[i] = m.Value()[i]
		}
	}
	vtc.Clock.Set(vec)
	vtc.Clock.Increment()
	return nil
}

func (vtc *VectorLocalClock) Snapshot() []uint64 {
	return vtc.Clock.Value()
}

func (lc *VectorLocalClock) SnapshotTS() DistributedClock[[]uint64] {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()
	return lc.Clock.Clone()
}
