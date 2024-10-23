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
	MachineGenerator    func(
		id uint64,
	) Machine[T, K, Message[T, K]]
	Machines []Machine[T, K, Message[T, K]]
}

func (factory *MachineFactory[T, K]) produce(assignedId uint64) Machine[T, K, Message[T, K]] {
	machine := factory.MachineGenerator(assignedId)
	machine.SetManager(
		timestamp.TsManagerNew[T, K, timestamp.DistributedClock[K], timestamp.LocalClock[K]](
			0, factory.LocalClockGenerator(),
		),
	)
	return machine
}

func (factory *MachineFactory[T, K]) InitAll() error {
	var i uint64 = 0
	factory.Machines = make([]Machine[T, K, Message[T, K]], factory.NumNode)

	for ; i < factory.NumNode; i++ {
		factory.Machines[i] = factory.produce(i)
		factory.Machines[i].Start()
	}

	return nil
}
