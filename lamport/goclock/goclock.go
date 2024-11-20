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
		) Machine[string, int] {
			return &MachineImpl[string, int]{
				data:         "None",
				listenCycle:  time.Millisecond * 300,
				nNodes:       numNode,
				nSubThread:   1,
				nPubThread:   1,
				pub_sub_size: pub_sub_size,
				id:           id,
			}
		},
	}

	if err := factory.InitAll(); err != nil {
		fmt.Println(err)
	}
	// fmt.Println("here??")
	ConfigSimpleDistributedStorage(&factory)
	factory.StartAll()
	TestSimpleDistributedStorage(&factory)
}

func ConfigSimpleDistributedStorage(factory *MachineFactory[string, int]) {
	mb := &MessageBroker{
		nExch:     factory.NumNode,
		exchanges: make(map[uint64]exchange.Exchange, factory.NumNode),
	}

	for _, m := range factory.Machines {
		exch, err := exchange.NewSingleBroker(0, 0, 0, 0, m.GetId())
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

func ConfigDelayDistributedStorage(factory *MachineFactory[string, int], delay_max time.Duration, delay_min time.Duration) {
	mb := &MessageBroker{
		nExch:     factory.NumNode,
		exchanges: make(map[uint64]exchange.Exchange, factory.NumNode),
	}

	for _, m := range factory.Machines {
		exch, err := exchange.NewSingleBroker(delay_max, delay_min, 0, 0, m.GetId())
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

func ConfigPoissonDelayDistributedStorage(factory *MachineFactory[string, int], avg_delay time.Duration) {
	mb := &MessageBroker{
		nExch:     factory.NumNode,
		exchanges: make(map[uint64]exchange.Exchange, factory.NumNode),
	}

	for _, m := range factory.Machines {
		exch, err := exchange.NewPoissonBroker(avg_delay, 0, 0, m.GetId())
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

func TestSimpleDistributedStorage(factory *MachineFactory[string, int]) {
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

func InitMachine(data string, listenCycle time.Duration, nNodes uint64, nSubThread uint32, nPubThread uint32, pub_sub_size uint32) func(uint64) Machine[string, int] {
	return func(
		id uint64,
	) Machine[string, int] {
		return &MachineImpl[string, int]{
			data:         data,
			listenCycle:  listenCycle,
			nNodes:       nNodes,
			nSubThread:   nSubThread,
			nPubThread:   nPubThread,
			pub_sub_size: pub_sub_size,
			id:           id,
		}
	}
}
