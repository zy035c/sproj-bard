package goclock

import (
	"lamport/timestamp"
	"time"
)

type MachineFactory[T Payload, K ClockDataType] struct {
	NumNode             uint64
	ListenCycle         time.Duration
	BufferSize          int
	LocalClockGenerator func() timestamp.LocalClock[K]
	MachineGenerator    func(id uint64) Machine[T, K]
	Machines            []Machine[T, K]
	Verbose             bool
}

func (factory *MachineFactory[T, K]) produce(assignedId uint64) Machine[T, K] {
	machine := factory.MachineGenerator(assignedId)
	machine.SetManager(
		timestamp.TsManagerNew[T, K, timestamp.DistributedClock[K], timestamp.LocalClock[K]](
			0, factory.LocalClockGenerator(),
		),
	)
	machine.SetVerbose(factory.Verbose)
	return machine
}

func (factory *MachineFactory[T, K]) InitAll() error {
	var i uint64 = 0
	factory.Machines = make([]Machine[T, K], factory.NumNode)

	for ; i < factory.NumNode; i++ {
		factory.Machines[i] = factory.produce(i)
	}

	return nil
}

func (factory *MachineFactory[T, K]) StartAll() {
	var i uint64 = 0
	for ; i < factory.NumNode; i++ {
		factory.Machines[i].Start()
	}
}
