package goclock

import "lamport/exchange"

type MessageBroker[M any] struct {
	exchanges map[uint64]exchange.Exchange[M]
	nExch     uint64
}

func (mb *MessageBroker[M]) Add(ex exchange.Exchange[M]) {
	mb.exchanges[ex.GetId()] = ex
}

func (mb *MessageBroker[M]) Get(id uint64) exchange.Exchange[M] {
	return mb.exchanges[id]
}

// func (mb *MessageBroker[M]) AsPub(m Machine[any, any, any]) exchange.Exchange[M] {
// 	return mb.exchanges[id]
// }

// func (mb *MessageBroker[M]) AsSub(id uint64) exchange.Exchange[M] {
// 	return mb.exchanges[id]
// }
