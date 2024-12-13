package goclock

import "lamport/exchange"

type MessageBroker struct {
	exchanges map[uint64]exchange.Exchange
	nExch     uint64
}

func (mb *MessageBroker) Add(ex exchange.Exchange) {
	mb.exchanges[ex.GetId()] = ex
}

func (mb *MessageBroker) Get(id uint64) exchange.Exchange {
	return mb.exchanges[id]
}

// func (mb *MessageBroker[M]) AsPub(m Machine[any, any, any]) exchange.Exchange[M] {
// 	return mb.exchanges[id]
// }

// func (mb *MessageBroker[M]) AsSub(id uint64) exchange.Exchange[M] {
// 	return mb.exchanges[id]
// }
