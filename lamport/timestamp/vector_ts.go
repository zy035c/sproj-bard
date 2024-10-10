package timestamp

import (
	"fmt"
	"lamport/utils"
	"sync"
)

type VectorTsCounter struct {
	vector []uint64
	mutex  sync.Mutex
	idx    uint32
	size   uint32
}

func New(sz uint32, idx uint32) (*VectorTsCounter, error) {
	if idx >= sz {
		return nil, fmt.Errorf("idx %d is eq or gt than size %d", idx, sz)
	}
	vector := make([]uint64, sz)
	return &VectorTsCounter{
		vector: vector,
		idx:    idx,
		size:   sz,
	}, nil
}

func (vtc *VectorTsCounter) forward() {
	vtc.mutex.Lock()
	defer vtc.mutex.Unlock()
	vtc.vector[vtc.idx]++
}

func (vtc *VectorTsCounter) adjust(m []uint64) error {
	vtc.mutex.Lock()
	defer vtc.mutex.Unlock()
	if len(m) != int(vtc.size) {
		return fmt.Errorf("ts has a size of %d, local ts has size %d", len(m), vtc.size)
	}

	for i, ts := range vtc.vector {
		if m[i] > ts {
			vtc.vector[i] = m[i]
		}
	}
	vtc.vector[vtc.idx]++
	return nil
}

func (vtc *VectorTsCounter) snapshot() []uint64 {
	return utils.SliceCpy[uint64](vtc.vector)
}

type TsOrder uint8

const (
	AFT TsOrder = iota
	BEF
	CON
)

type Version[T any] struct {
	data      T
	timestamp []uint64
}

type TsManager[T any] struct {
	VerChain []Version[T]
}
