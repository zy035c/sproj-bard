package simulation

import (
	"fmt"
	"lamport/chart"
	"lamport/goclock"
	"lamport/timestamp"
	"lamport/utils"
	"math/rand"
	"sync/atomic"
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
	PossionRandomSimulation()
	// goclock.Main()
}

type Machine[T goclock.Payload, K goclock.ClockDataType] goclock.Machine[T, K]

func Config() *goclock.MachineFactory[string, int] {
	// full connected network
	var numNode uint64 = 12
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

func PossionRandomSimulation() {
	fmt.Printf("----- PossionRandomSimulation -----\n")
	factory := Config()
	if err := factory.InitAll(); err != nil {
		fmt.Println(err)
	}
	goclock.ConfigSimpleDistributedStorage(factory)
	factory.StartAll()

	var counter atomic.Int64
	counter.Store(1)

	var acc_time float64 = 0

	for {
		interval := conf.PoissonInterval()
		acc_time += float64(interval)
		if interval > 0 {
			fmt.Printf("Sleeping for %v\n", interval)
			time.Sleep(interval)
		}

		epoch := int(counter.Load())
		event := Event{
			Epoch: epoch,
			Mid:   rand.Intn(int(factory.NumNode)),
		}

		if utils.RandomFloat32(0.0, 1.0) < conf.ReadWriteRatio {
			event.Etype = READ
			event.Op = func() {
				fmt.Printf("--- Simulating: READ at Machine%v, Epoch%v\n", event.Mid, epoch)
				PlotFlawMetric(epoch, factory.Machines, acc_time, dyn_chart)
				RandomSampleNodeVersionChain(factory.Machines, event.Mid)
				PrintCycleMetric(epoch, factory.Machines)
			}
		} else {
			event.Etype = WRITE
			event.Op = func() {
				fmt.Printf("--- Simulating: WRITE at Machine%v, Epoch%v\n", event.Mid, epoch)
				factory.Machines[event.Mid].LocalEvent(func(data string) string {
					return fmt.Sprintf("Epoch%v", epoch)
				})
			}
			counter.Add(1)
		}
		event.Op()
	}
}

func FixedIntervalSimulation() {
	// interval := conf.AvgInterval
}

func PossionRandomSimulationWithDelay() {
	fmt.Printf("----- PossionRandomSimulationWithDelay -----\n")
	// TODO
}
