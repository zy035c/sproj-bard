package goclock

import (
	"lamport/timestamp"
	"time"
)

type MachineFactory[T Payload, K ClockDataType] struct {
	NumNode             uint64
	Data                T
	ListenCycle         time.Duration
	BufferSize          int
	LocalClockGenerator func() timestamp.LocalClock[K]
	Machines            []Machine[T, K]
	Channels            []chan Message[T, K]
}

func (factory *MachineFactory[T, K]) produce(assignedId uint64) (Machine[T, K], chan Message[T, K]) {
	recv := make(chan Message[T, K], factory.BufferSize)
	machine := MachineImpl[T, K, ClockAbbr[K], timestamp.LocalClock[K]]{
		data:        factory.Data,
		listenCycle: factory.ListenCycle,
		recv:        recv,
		id:          assignedId,
		nNodes:      factory.NumNode,
		manager: timestamp.TsManagerNew[T, K, ClockAbbr[K], timestamp.LocalClock[K]](
			0, factory.LocalClockGenerator(),
		),
	}

	return &machine, recv
}

func (factory *MachineFactory[T, K]) StartAll() error {
	var i uint64 = 0
	factory.Channels = make([]chan Message[T, K], factory.NumNode)
	factory.Machines = make([]Machine[T, K], factory.NumNode)

	for ; i < factory.NumNode; i++ {
		factory.Machines[i], factory.Channels[i] = factory.produce(i)
	}

	for i, machine := range factory.Machines {
		tmp := append(
			make([]chan Message[T, K], 0, factory.NumNode-1),
			factory.Channels[:i]...,
		)
		tmp = append(tmp, factory.Channels[i+1:]...)
		if err := machine.SetSend(tmp); err != nil {
			return err
		}
		machine.Start()
	}

	return nil
}
