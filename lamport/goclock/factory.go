package goclock

import (
	"fmt"
	"lamport/timestamp"
	"lamport/utils"
	"reflect"
	"sync"
	"time"
)

type MachineFactory[T any, ClockType timestamp.DistributedClock[T]] struct {
	Clock ClockType
}

type trueMachineFac[T Payload, K any, U timestamp.DistributedClock[K], G timestamp.LocalClock[K]] struct {
	manager *timestamp.TsManager[T, K, U, G]
}

func (factory MachineFactory[T, ClockType]) UseDataType(dataType PayloadType) *trueMachineFac[string, T, ClockType, timestamp.LocalClock[T]] {
	switch dataType {
	case StringType:
		return &trueMachineFac[string, T, ClockType, timestamp.LocalClock[T]]{}
	}

	return nil
}

type MachineImpl[T Payload, K any, U timestamp.DistributedClock[K], G timestamp.LocalClock[K]] struct {
	data        T // possession
	mutex       sync.Mutex
	listenCycle time.Duration
	channels    []chan Message[K, T]
	id          uint64
	nNodes      uint64
	manager     *timestamp.TsManager[T, K, U, G]
}

func (m *MachineImpl[T, K, U, G]) getData() T {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.data
}

func (m *MachineImpl[T, K, U, G]) Start() {
	//
	go func() {

	}()

	go func() { // listen loop
		for {

		}
	}()
}

func (m *MachineImpl[T, K, U, G]) LocalEvent(event func(data T) T) {
	m.data = event(m.data)
}

func (m *MachineImpl[T, K, U, G]) Listen() {
	interval := m.listenCycle
	channels := m.channels

	for {
		cases := make([]reflect.SelectCase, m.nNodes)
		for i, ch := range channels {
			cases[i] = reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(ch),
			}
		}

		// add timeout case
		cases[m.nNodes] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(time.After(interval)),
		}

		chosen, value, ok := reflect.Select(cases)
		if chosen < int(m.nNodes) && ok {
			m.handleMsg(utils.ReflectConvert[Message[T, K]](value))
		} else if chosen == int(m.nNodes) {
			fmt.Println("Timeout, restarting loop...")
		}
	}
}

func (m *MachineImpl[T, K, U, G]) handleMsg(msg *Message[T, K]) {
	if msg == nil {
		return
	}
	// m.manager.Add(timestamp.Version[T, K, U]{
	// 	*msg,
	// })
}

func (factory *trueMachineFac[T, K, U, G]) produce() Machine {

}
