package timestamp

import "sync"

type KeyValPair struct {
	Key string
	Val string
}

type VectorTsManager struct {
	*TsManager[KeyValPair, []uint64, *VectorClock, *VectorLocalClock]
}

type LamportTsManager struct {
	*TsManager[KeyValPair, int, *LamportClock, *LamportLocalClock]
}

func Main() {
	testSimpleLamport()
	testSimpleVector()
}

func testSimpleLamport() {

	AddVersion := func(key string, value string, lc int, mngr *LamportTsManager) {
		mngr.Add(
			Version[KeyValPair, int, *LamportClock]{
				data: KeyValPair{
					Key: key,
					Val: value,
				},
				timestamp: &LamportClock{
					counter: lc,
				},
			})
	}

	manager := &LamportTsManager{TsManagerNew[KeyValPair, int, *LamportClock, *LamportLocalClock](
		1024, &LamportLocalClock{
			&LamportClock{
				counter: 0,
			},
			sync.Mutex{},
		})}

	AddVersion("Greeting", "Hello", 4, manager)
	AddVersion("Greeting", "Zdrazdvui", 2, manager)
	AddVersion("Greeting", "Nihao", 1, manager)

	manager.PrintVersionChain()
}

func testSimpleVector() {

	mngr := &VectorTsManager{TsManagerNew[KeyValPair, []uint64, *VectorClock, *VectorLocalClock](
		1024, &VectorLocalClock{
			&VectorClock{
				vector: []uint64{0, 0, 0, 0, 0},
				idx:    2,
				size:   5,
			},
			sync.Mutex{},
			2,
			5,
		})}

	AddVersion := func(key string, value string, lc []uint64, mngr *VectorTsManager) {
		mngr.Add(
			Version[KeyValPair, []uint64, *VectorClock]{
				data: KeyValPair{
					Key: key,
					Val: value,
				},
				timestamp: &VectorClock{
					vector: lc,
					idx:    2,
					size:   5,
				},
			})
	}

	AddVersion("Greeting", "Zdrazdvui", []uint64{2, 2, 4, 4, 5}, mngr)
	AddVersion("Greeting", "Nihao", []uint64{43, 2, 8, 78, 4}, mngr)
	AddVersion("Greeting", "Hello", []uint64{1, 2, 3, 4, 1}, mngr)

	mngr.PrintVersionChain()
}
