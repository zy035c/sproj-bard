package goclock

import (
	"encoding/json"
	"fmt"
	"lamport/exchange"
	"lamport/option"
	"lamport/timestamp"
	"lamport/utils"
	"sync"
	"time"
)

type VersionPtr[T any, K any] *timestamp.Version[T, K]

type MachineImpl[T Payload, K ClockDataType] struct { // in fact M should be bound by Version[T, K]
	data        T // possession
	mutex       sync.Mutex
	listenCycle time.Duration
	sub         []exchange.Exchange
	pub         []exchange.Exchange
	id          uint64
	nNodes      uint64
	manager     *timestamp.TsManager[T, K, timestamp.DistributedClock[K], timestamp.LocalClock[K]]

	nSubThread   uint32
	nPubThread   uint32
	pub_sub_size uint32
	pubPool      utils.ThreadPool

	verbose bool
}

func (m *MachineImpl[T, K]) GetData() T {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.data
}

func (m *MachineImpl[T, K]) Init() {
	// m.pub = make([]exchange.Exchange)
	// m.sub = make([]exchange.Exchange, m.pub_sub_size)
}

func (m *MachineImpl[T, K]) Start() {
	if m.nSubThread == 0 {
		m.nSubThread = 1
	}
	if m.nPubThread == 0 {
		m.nPubThread = 1
	}
	m.pubPool.Init(int(m.nPubThread))
	m.Listen()
}

func (m *MachineImpl[T, K]) Stop() {
	// state machine (e.g. simulated down/offline)
	m.mutex.Lock()
	defer m.mutex.Unlock()
}

func (m *MachineImpl[T, K]) LocalEvent(event func(data T) T) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data = event(m.data) // ? TODO
	cur_version := timestamp.NewVersion[T, K](
		m.data,
		m.manager.LocalClk.SnapshotTS(),
		m.id,
	)
	m.manager.Add(cur_version.Clone())
	m.Broadcast(&cur_version)
}

func (m *MachineImpl[T, K]) Broadcast(msg VersionPtr[T, K]) {
	for _, pub := range m.pub {
		v := (*msg).Clone()
		m.publish(
			&v,
			pub,
		)
	}
}

func (m *MachineImpl[T, K]) publish(msg VersionPtr[T, K], ex exchange.Exchange) {
	m.pubPool.Submit(func() any {
		bytestream, err := json.Marshal(*msg)
		if err != nil {
			panic(err)
		}
		res := ex.Put(bytestream)
		if m.verbose {
			fmt.Printf("+ Machine %v Pub msg %v to exch %v res=%v\n", m.id, *msg, ex.GetId(), res)
		}
		return res
	})
}

func (m *MachineImpl[T, K]) poll() *option.Option[VersionPtr[T, K]] {
	for _, ch := range m.sub {
		select {
		case bytestream := <-ch.C():
			ver := &timestamp.Version[T, K]{}
			json.Unmarshal(bytestream, ver)
			return option.Some[VersionPtr[T, K]](ver)
		default:
			//
		}
	}
	return option.None[VersionPtr[T, K]]()
}

func (m *MachineImpl[T, K]) Listen() {
	for i := 0; i < int(m.nSubThread); i++ {
		go func() {
			for {
				res := m.poll()
				if res.Has() {
					msg := res.Unwrap()
					if m.verbose {
						fmt.Printf("+ Machine %v Sub msg %v\n", m.id, msg)
					}
					m.handleMsg(msg)
				} else {
					time.Sleep(m.listenCycle)
				}
			}
		}()
	}
}

func (m *MachineImpl[T, K]) handleMsg(msg VersionPtr[T, K]) {
	// if !ok {
	// 	panic("")
	// }
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.manager.Add(*msg)
}

func (m *MachineImpl[T, K]) BindSub(recv exchange.Exchange) {
	// m.mutex.Lock()
	// defer m.mutex.Unlock()
	m.sub = append(m.sub, recv)
}

func (m *MachineImpl[T, K]) BindPub(send exchange.Exchange) error {
	// m.mutex.Lock()
	// defer m.mutex.Unlock()
	m.pub = append(m.pub, send)
	return nil
}

func (m *MachineImpl[T, K]) PrintInfo() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	fmt.Printf("Machine %v\n- Local Clock %v\n", m.id, m.manager.LocalClk.SnapshotTS())
	if m.verbose {
		fmt.Printf("- Version Chain\n")
		m.manager.VerChain.Traverse()
		fmt.Printf("Machine %v Ends\n\n", m.id)
	}
}

func (m *MachineImpl[T, K]) GetId() uint64 {
	return m.id
}

func (m *MachineImpl[T, K]) SetManager(mngr *timestamp.TsManager[T, K, timestamp.DistributedClock[K], timestamp.LocalClock[K]]) {
	m.manager = mngr
}

func (m *MachineImpl[T, K]) GetVersionChainData() []T {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.manager.GetVersionChainData()
}

func (m *MachineImpl[T, K]) SetVerbose(v bool) {
	m.verbose = v
}

var _ Machine[string, int] = &MachineImpl[string, int]{}
