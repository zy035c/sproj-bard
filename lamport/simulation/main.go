package simulation

import (
	"fmt"
	"lamport/chart"
	"lamport/goclock"
	"lamport/timestamp"
	"lamport/utils"
	"math/rand"
	"time"
)

var dyn_chart = chart.Init()

var conf = SimConfig{
	ReadWriteRatio: 0.8,
	AvgInterval:    500 * time.Millisecond,
	AvgDelay:       60 * time.Millisecond,
}

func Main() {
	go dyn_chart.Main()
	PossionRandomSimulationWithDelay()
	// goclock.Main()
}

type Machine[T goclock.Payload, K goclock.ClockDataType] goclock.Machine[T, K]

func Config() *goclock.MachineFactory[string, int] {
	// full connected network
	var numNode uint64 = 5
	var pub_sub_size uint32 = 1024

	factory := goclock.MachineFactory[string, int]{
		NumNode: numNode,
		LocalClockGenerator: func() timestamp.LocalClock[int] {
			return &timestamp.LamportLocalClock{
				Clock: &timestamp.LamportClock{},
			}
		},
		MachineGenerator: goclock.InitMachine("None", time.Millisecond*300, numNode, 1, 10, pub_sub_size),
		Verbose:          false,
	}

	return &factory
}

func Task(factory *goclock.MachineFactory[string, int], epoch int) {
	if epoch%10 == 0 {
		fmt.Printf("\n\n********** Epoch %v PrintInfo **********\n\n", epoch)
		for _, m := range factory.Machines {
			m.PrintInfo()
			fmt.Printf("- Score=%v\n", CalcScore(epoch, m))
			fmt.Printf("\n")
		}
		fmt.Printf("\n********** PrintInfo End **********\n\n\n")
	}
	curId := rand.Intn(int(factory.NumNode))
	factory.Machines[curId].LocalEvent(func(data string) string {
		return fmt.Sprintf("Epoch%v", epoch)
	})
	fmt.Printf("- Simulating: Local Event at Machine%v, Epoch%v\n", curId, epoch)
}

func PossionRandomSimulation() {
	fmt.Printf("----- PossionRandomSimulation -----\n")
	factory := Config()
	if err := factory.InitAll(); err != nil {
		fmt.Println(err)
	}
	goclock.ConfigSimpleDistributedStorage(factory)
	factory.StartAll()

	epoch := 1

	var acc_time float64 = 0

	for {
		interval := conf.PoissonInterval()
		acc_time += float64(interval)
		if interval > 0 {
			fmt.Printf("Sleeping for %v\n", interval)
			time.Sleep(interval)
		}
		if utils.RandomFloat32(0.0, 1.0) < conf.ReadWriteRatio {
			PlotScore(epoch, factory.Machines, acc_time, dyn_chart)
		} else {
			Task(factory, epoch)
			epoch++
		}
	}
}

func FixedIntervalSimulation() {
	fmt.Printf("----- FixedIntervalSimulation -----\n")
	factory := Config()

	factory.InitAll()
	goclock.ConfigSimpleDistributedStorage(factory)
	factory.StartAll()

	var acc_time float64 = 0

	epoch := 1

	for {
		interval := conf.AvgInterval
		if interval > 0 {
			fmt.Printf("Sleeping for %v\n", interval)
			time.Sleep(conf.AvgInterval)
		}
		if utils.RandomFloat32(0.0, 1.0) < conf.ReadWriteRatio {
			PlotScore(epoch, factory.Machines, acc_time, dyn_chart)
		} else {
			Task(factory, epoch)
			epoch++
		}
	}
}

func PossionRandomSimulationWithDelay() {
	fmt.Printf("----- PossionRandomSimulationWithDelay -----\n")
	factory := Config()
	if err := factory.InitAll(); err != nil {
		fmt.Println(err)
	}
	goclock.ConfigPoissonDelayDistributedStorage(factory, conf.AvgDelay)
	factory.StartAll()

	epoch := 1

	var acc_time float64 = 0

	for {
		interval := conf.PoissonInterval()
		acc_time += float64(interval)
		if interval > 0 {
			fmt.Printf("Sleeping for %v\n", interval)
			time.Sleep(interval)
		}
		if utils.RandomFloat32(0.0, 1.0) < conf.ReadWriteRatio {
			PlotScore(epoch, factory.Machines, acc_time, dyn_chart)
		} else {
			Task(factory, epoch)
			epoch++
		}
	}
}
