package timestamp

import (
	"lamport/utils"
)

type TsManager[T any, K any, U DistributedClock[K]] struct {
	VerChain utils.LinkList[Version[T, K, U]]
	size     uint32
	// TsCmpFunc[T]
}

func TsManagerNew[T any, K any, U DistributedClock[K]](sz int) *TsManager[T, K, U] {
	// var cmp TsCmpFunc[T]
	// if ts_cmp_func == nil {
	// 	cmp = defaultCmp
	// } else {
	// 	cmp = *ts_cmp_func
	// }

	return &TsManager[T, K, U]{
		VerChain: *utils.NewLinkList[Version[T, K, U]](),
		size:     0,
	}
}

func (tsm *TsManager[T, K, U]) Add(m Version[T, K, U]) {

	tsm.VerChain.InsertBefore(&m,
		func(v1, v2 *Version[T, K, U]) bool {
			// fmt.Printf("Comparing %v and %v\n", v1, v2)
			order := v1.timestamp.DefaultTsCmp(v2.timestamp)
			return order == CON || order == AFT
		})

	tsm.size++
}

func (tsm TsManager[T, K, U]) PrintVersionChain() {
	tsm.VerChain.Traverse()
}
