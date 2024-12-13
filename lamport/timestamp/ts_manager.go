package timestamp

import (
	"fmt"
	"lamport/utils"
)

type TsManager[T any, K any, U DistributedClock[K], G LocalClock[K]] struct {
	VerChain *utils.LinkList[Version[T, K]]
	LocalClk G
	size     uint32
}

func TsManagerNew[T any, K any, U DistributedClock[K], G LocalClock[K]](
	sz int,
	local_clock G,
) *TsManager[T, K, U, G] {
	return &TsManager[T, K, U, G]{
		VerChain: utils.NewLinkList[Version[T, K]](),
		size:     0,
		LocalClk: local_clock,
	}
}

func (tsm *TsManager[T, K, U, G]) Add(m Version[T, K]) {

	tsm.LocalClk.Adjust(m.GetTs())

	tsm.VerChain.InsertBefore(&m,
		func(v1, v2 *Version[T, K]) bool {
			// fmt.Printf("Comparing %v and %v\n", v1, v2)
			order := v1.timestamp.DefaultTsCmp(v2.timestamp)
			return order == CON || order == AFT
		})

	tsm.size++
}

func (tsm TsManager[T, K, U, G]) PrintVersionChain() {
	fmt.Printf("Current ts - %v\n", tsm.LocalClk.Snapshot())
	tsm.VerChain.Traverse()
}

func (tsm TsManager[T, K, U, G]) GetVersionChainData() []T {
	li := []T{}
	for _, x := range tsm.VerChain.ToList() {
		li = append(li, x.GetData())
	}
	return li
}
