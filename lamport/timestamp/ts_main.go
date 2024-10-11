package timestamp

type KeyValPair struct {
	Key string
	Val string
}

func Main() {
	testSimpleLamport()
	testSimpleVector()
}

func testSimpleLamport() {

	AddVersion := func(key string, value string, lc int, mngr *TsManager[KeyValPair, int, LamportClock]) {
		mngr.Add(
			Version[KeyValPair, int, LamportClock]{
				data: &KeyValPair{
					Key: key,
					Val: value,
				},
				timestamp: LamportClock{
					counter: lc,
				},
			})
	}

	mngr := TsManagerNew[KeyValPair, int, LamportClock](1024)

	AddVersion("Greeting", "Hello", 4, mngr)
	AddVersion("Greeting", "Zdrazdvui", 2, mngr)
	AddVersion("Greeting", "Nihao", 1, mngr)

	mngr.PrintVersionChain()
}

func testSimpleVector() {
	mngr := TsManagerNew[KeyValPair, []uint64, VectorClock](1024)

	AddVersion := func(key string, value string, lc []uint64, mngr *TsManager[KeyValPair, []uint64, VectorClock]) {
		mngr.Add(
			Version[KeyValPair, []uint64, VectorClock]{
				data: &KeyValPair{
					Key: key,
					Val: value,
				},
				timestamp: VectorClock{
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
