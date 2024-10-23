package goclock

import (
	"fmt"
	"lamport/exchange"
	"lamport/option"
	"lamport/timestamp"
	"sync"
	"time"
)

type MachineImpl[T Payload, K ClockDataType, U timestamp.DistributedClock[K], G timestamp.LocalClock[K]] struct { // in fact M should be bound by Message[T, K]
	data        T // possession
	mutex       sync.Mutex
	listenCycle time.Duration
	sub         []exchange.Exchange[Message[T, K]]
	pub         []exchange.Exchange[Message[T, K]]
	id          uint64
	nNodes      uint64
	manager     *timestamp.TsManager[T, K, timestamp.DistributedClock[K], timestamp.LocalClock[K]]

	nSubThread   uint32
	nPubThread   uint32
	pub_sub_size uint32
	// pubPool    utils.ThreadPool
}

func (m *MachineImpl[T, K, U, G]) GetData() T {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.data
}

func (m *MachineImpl[T, K, U, G]) Start() {
	if m.nSubThread == 0 {
		m.nSubThread = 1
	}
	if m.nPubThread == 0 {
		m.nPubThread = 1
	}
	m.pub = make([]exchange.Exchange[Message[T, K]], m.pub_sub_size)
	m.sub = make([]exchange.Exchange[Message[T, K]], m.pub_sub_size)
	// m.pubPool.Init(int(m.nPubThread))
	m.Listen()
}

func (m *MachineImpl[T, K, U, G]) Stop() {
	// state machine (e.g. simulated down/offline)
}

func (m *MachineImpl[T, K, U, G]) LocalEvent(event func(data T) T) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data = event(m.data)
	cur_version := timestamp.NewVersion[T, K](
		m.data,
		m.manager.LocalClk.SnapshotTS().(U),
		m.id,
	)
	m.manager.Add(cur_version.Clone())
	for _, pub := range m.pub {
		// m.publish(
		// 	cur_version.Clone(),
		// 	pub,
		// )
		cl := cur_version.Clone()
		res := pub.Put(cl)
		fmt.Printf("+ Pub msg %v to exch %v res=%v\n", cur_version, pub.GetId(), res)
	}
	m.Broadcast(cur_version)
}

func (m *MachineImpl[T, K, U, G]) Broadcast(msg Message[T, K]) {
	// for _, pub := range m.pub {
	// 	m.publish(
	// 		msg.(timestamp.Version[T, K]).Clone(),
	// 		pub,
	// 	)
	// }
}

func (m *MachineImpl[T, K, U, G]) publish(msg Message[T, K], ex exchange.Exchange[Message[T, K]]) {
	// m.pubPool.Submit(func() any {
	// 	res := ex.Put(msg.(M))
	// 	fmt.Println("+ Pub msg %v to exch &v res=%v\n", msg, ex.GetId(), res)
	// 	return res
	// }).Then(func(_ any) any {
	// 	// TODO
	// 	return nil
	// })
	// res := ex.Put(msg.(M))
	// fmt.Printf("+ Pub msg %v to exch %v res=%v\n", msg, ex.GetId(), res)
}

func (m *MachineImpl[T, K, U, G]) poll() *option.Option[Message[T, K]] {
	// var cases []reflect.SelectCase
	// for _, ch := range m.sub {
	// 	cases = append(cases, reflect.SelectCase{
	// 		Dir:  reflect.SelectRecv,
	// 		Chan: reflect.ValueOf(ch),
	// 	})
	// }
	// cases = append(cases, reflect.SelectCase{
	// 	Dir: reflect.SelectDefault,
	// })
	// chosen, recv, ok := reflect.Select(cases)
	// if chosen < len(m.sub) && ok {
	// 	return utils.ReflectConvert[Message[T, K]](recv)
	// }
	for _, ch := range m.sub {
		select {
		case elem := <-ch.C():
			return option.Some[Message[T, K]](elem)
		default:
			//
		}
	}
	return option.None[Message[T, K]]()
}

func (m *MachineImpl[T, K, U, G]) Listen() {
	for i := 0; i < int(m.nSubThread); i++ {
		go func() {
			for {
				res := m.poll()
				if res.Has() {
					msg := res.Unwrap()
					fmt.Printf("+ Machine %v Receiving msg %v\n", m.id, msg)
					m.handleMsg(msg)
				} else {
					time.Sleep(m.listenCycle)
				}
			}
		}()
	}
}

func (m *MachineImpl[T, K, U, G]) handleMsg(msg Message[T, K]) {
	v, ok := msg.(timestamp.Version[T, K])
	if !ok {
		panic("")
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.manager.Add(v)
}

func (m *MachineImpl[T, K, U, G]) BindSub(recv exchange.Exchange[Message[T, K]]) {
	m.sub = append(m.sub, recv)
}

func (m *MachineImpl[T, K, U, G]) BindPub(send exchange.Exchange[Message[T, K]]) error {
	// if len(send) != int(m.nNodes)-1 {
	// 	return fmt.Errorf("! SetSend []chan size %v is not valid", len(send))
	// }
	m.pub = append(m.pub, send)
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

func (m *MachineImpl[T, K, U, G]) GetId() uint64 {
	return m.id
}

func (m *MachineImpl[T, K, U, G]) SetManager(mngr *timestamp.TsManager[T, K, timestamp.DistributedClock[K], timestamp.LocalClock[K]]) {
	m.manager = mngr
}
