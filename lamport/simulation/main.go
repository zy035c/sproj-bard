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
	ReadWriteRatio: 0.6,
	AvgInterval:    60 * time.Millisecond,
	AvgDelay:       40 * time.Millisecond,
	PLR:            0.00388,
}

func Main() {
	go dyn_chart.Main()
	PossionRandomSimulation()
	// goclock.Main()
}

type Machine[T goclock.Payload, K goclock.ClockDataType] goclock.Machine[T, K]

func Config() *goclock.MachineFactory[string, int] {
	// full connected network
	var numNode uint64 = 6
	var pub_sub_size uint32 = 1024

	factory := goclock.MachineFactory[string, int]{
		NumNode: numNode,
		LocalClockGenerator: func() timestamp.LocalClock[int] {
			return &timestamp.LamportLocalClock{
				Clock: &timestamp.LamportClock{},
			}
		},
		MachineGenerator: goclock.InitMachine("None", time.Millisecond*0, numNode, 0, 0, pub_sub_size),
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
	goclock.ConfigPoissonDelayDistributedStorage(factory, conf.AvgDelay, conf.PLR)
	factory.StartAll()
	InitMetrics()

	var counter atomic.Int64
	counter.Store(1)

	hist := &ClientHistory{}

	var acc_time float64 = 0
	real_epoch := 1

	for {
		interval := conf.PoissonInterval()
		acc_time += float64(interval)
		if interval > 0 {
			time.Sleep(interval)
		}

		event := Event{
			Epoch: real_epoch,
			Vid:   int(counter.Load()),
			Mid:   rand.Intn(int(factory.NumNode)),
		}

		var version_chain []string

		if utils.RandomFloat32(0.0, 1.0) < conf.ReadWriteRatio {
			event.Etype = READ
			event.Op = func() {
				version_chain = factory.Machines[event.Mid].GetVersionChainData()
				// fmt.Printf("--- Simulating: READ at Machine%v, Epoch%v\n", event.Mid, epoch)

				VerChainStrToInt(version_chain)
				// fmt.Println("- Sampling Node", event.Mid, strings.Join(version_chain, "->"))
			}
		} else {
			event.Etype = WRITE
			vid := int(counter.Load())
			counter.Add(1)
			event.Op = func() {
				// fmt.Printf("--- Simulating: WRITE at Machine%v, Epoch%v\n", event.Mid, epoch)
				factory.Machines[event.Mid].LocalEvent(func(data string) string {
					return fmt.Sprintf("Vid%v", vid)
				})
			}
			hist.Add(event.Mid, &event)
		}
		event.Op()

		lastWriteVid := -1
		lastWrite := hist.GetLast(event.Mid)
		if lastWrite != nil {
			lastWriteVid = lastWrite.Vid
		}

		version_chain = factory.Machines[event.Mid].GetVersionChainData()
		vid := int(counter.Load())

		go UpdateMetric(version_chain, vid, lastWriteVid, event.Mid, real_epoch)
		real_epoch += 1
	}
}

func FixedIntervalSimulation() {
	// interval := conf.AvgInterval
}

func UpdateMetric(version_chain []string, vid int, lastWriteVid int, mid int, iter int) {
	a, b, c := PrintMetrics(vid, VerChainStrToInt(version_chain), lastWriteVid, mid)
	if iter%1000 == 0 {
		fmt.Println("-----", iter, "-----")
		// PlotFlawMetric(int(counter.Load()), factory.Machines, acc_time, dyn_chart)
		fmt.Println("~MRC Metric", a)
		fmt.Println("~RYWC Metric", b)
		fmt.Println("~MRW Metric", c)
		fmt.Println("-----", iter, "-----")
	}
}

func PossionRandomSimulationWithDelay() {
	fmt.Printf("----- PossionRandomSimulationWithDelay -----\n")
	// TODO
}
