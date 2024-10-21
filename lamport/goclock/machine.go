package goclock

import (
	"fmt"
	"lamport/timestamp"
	"lamport/utils"
	"sync"
	"time"
)

type MachineImpl[T Payload, K any, U timestamp.DistributedClock[K], G timestamp.LocalClock[K]] struct {
	data        T // possession
	mutex       sync.Mutex
	listenCycle time.Duration
	recv        chan Message[T, K]
	send        []chan Message[T, K]
	id          uint64
	nNodes      uint64
	manager     *timestamp.TsManager[T, K, U, G]
}

func (m *MachineImpl[T, K, U, G]) GetData() T {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.data
}

func (m *MachineImpl[T, K, U, G]) Start() {
	go m.Listen()
}

func (m *MachineImpl[T, K, U, G]) Stop() {
	//
}

func (m *MachineImpl[T, K, U, G]) LocalEvent(event func(data T) T) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data = event(m.data)
	cur_version := timestamp.NewVersion[T, K, U](
		m.data,
		m.manager.LocalClk.SnapshotTS().(U),
		m.id,
	)
	m.manager.Add(*cur_version)
	m.Broadcast(cur_version)
}

func (m *MachineImpl[T, K, U, G]) Listen() {
	interval := m.listenCycle
	recv := m.recv
	timer := time.NewTimer(0)

	for {
		timer.Reset(interval)
		select {
		case msg := <-recv:
			m.handleMsg(msg)
		case <-timer.C:
			// pass
		}
	}
}

func (m *MachineImpl[T, K, U, G]) handleMsg(msg Message[T, K]) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	v := msg.(timestamp.Version[T, K, U])
	fmt.Printf("+ Machine %v Receiving msg %v\n", m.id, v)
	m.manager.Add(v)
}

func (m *MachineImpl[T, K, U, G]) Broadcast(msg Message[T, K]) {
	for i, sendch := range m.send {
		sendch <- *msg.(utils.Cloneable[*timestamp.Version[T, K, ClockAbbr[K]]]).Clone()
		chanId := i
		if chanId >= int(m.id) {
			chanId++
		}
		fmt.Printf("+ Machine %v Broadcasting msg to channel #%v\n", m.id, chanId)
	}
}

func (m *MachineImpl[T, K, U, G]) SetSend(send []chan Message[T, K]) error {
	// if len(send) != int(m.nNodes)-1 {
	// 	return fmt.Errorf("! SetSend []chan size %v is not valid", len(send))
	// }
	m.send = send
	return nil
}

func (m *MachineImpl[T, K, U, G]) PrintInfo() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	fmt.Printf("Machine %v\n- Local Clock %v\n- Version Chain\n", m.id, m.manager.LocalClk.SnapshotTS())
	// fmt.Printf("- Data %v", m.GetData())
	m.manager.VerChain.Traverse()
	fmt.Printf("Machine %v Ends\n\n", m.id)
}
