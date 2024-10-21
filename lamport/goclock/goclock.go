package goclock

import (
	"fmt"
	"lamport/timestamp"
	"math/rand"
	"time"
)

func Main() {
	factory := MachineFactory[string, int]{
		NumNode:     5,
		Data:        "Default",
		ListenCycle: time.Millisecond * 500,
		BufferSize:  1024,
		LocalClockGenerator: func() timestamp.LocalClock[int] {
			return &timestamp.LamportLocalClock{
				Clock: &timestamp.LamportClock{},
			}
		},
	}

	if err := factory.StartAll(); err != nil {
		fmt.Println(err)
	}
	SimulateA(&factory)
}

func SimulateA(factory *MachineFactory[string, int]) {

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
