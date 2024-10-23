package goclock

import (
	"fmt"
	"lamport/exchange"
	"lamport/timestamp"
	"math/rand"
	"time"
)

func Main() {

	var numNode uint64 = 5
	var pub_sub_size uint32 = 1024

	factory := MachineFactory[string, int]{
		NumNode: numNode,
		LocalClockGenerator: func() timestamp.LocalClock[int] {
			return &timestamp.LamportLocalClock{
				Clock: &timestamp.LamportClock{},
			}
		},
		MachineGenerator: func(
			id uint64,
		) Machine[string, int, Message[string, int]] {
			return &MachineImpl[string, int, *timestamp.LamportClock, *timestamp.LamportLocalClock]{
				data:         "None",
				listenCycle:  time.Millisecond * 300,
				nNodes:       numNode,
				nSubThread:   1,
				nPubThread:   1,
				pub_sub_size: pub_sub_size,
			}
		},
	}

	if err := factory.InitAll(); err != nil {
		fmt.Println(err)
	}
	ConfigSimpleDistributedStorage(&factory)
	SimulateSimpleDistributedStorage(&factory)
}

func ConfigSimpleDistributedStorage(factory *MachineFactory[string, int]) {
	mb := &MessageBroker[timestamp.Version[string, int]]{
		nExch:     factory.NumNode,
		exchanges: make(map[uint64]exchange.Exchange[timestamp.Version[string, int]]),
	}

	for _, m := range factory.Machines {
		exch, err := exchange.NewExchangeImpl[timestamp.Version[string, int]](0, 0, 0, 0, m.GetId())
		if err != nil {
			panic("Failed Creating ExchangeImpl")
		}
		mb.Add(exch)
		m.BindSub(exch)
	}

	for _, m := range factory.Machines {
		for k, v := range mb.exchanges {
			if k != m.GetId() {
				m.BindPub(v)
			}
		}
	}
}

func SimulateSimpleDistributedStorage(factory *MachineFactory[string, int]) {

	fmt.Printf("----- SimulateSimpleDistributedStorage -----\n")
	round := 0

	for {

		if round%5 == 0 {
			fmt.Printf("----- Epoch %v -----\n", round)
			for _, m := range factory.Machines {
				m.PrintInfo()
			}
			fmt.Printf("----- Epoch %v End-----\n", round)
		}

		curId := rand.Intn(int(factory.NumNode))

		factory.Machines[curId].LocalEvent(func(data string) string {
			return fmt.Sprintf("Round%v", round)
		})

		fmt.Printf("- Simulating: Local Event at Machine%v, Round%v\n", curId, round)

		round++
		time.Sleep(time.Second * 2)
	}
}
