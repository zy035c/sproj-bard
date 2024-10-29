package goclock

import (
	"lamport/exchange"
	"lamport/option"
	"lamport/timestamp"
)

type PayloadType uint8

type KeyValPair struct {
	Key string
	Val string
}

const (
	StringType PayloadType = iota
	FunctionType
	KVPairType
)

type Payload interface {
	string | KeyValPair | func() any
}

type ClockDataType interface {
	int | []uint64
}

type Machine[T Payload, K ClockDataType] interface {
	Start()
	Stop()
	Init()

	publish(VersionPtr[T, K], exchange.Exchange)
	poll() *option.Option[VersionPtr[T, K]]
	Listen()
	Broadcast(VersionPtr[T, K])

	BindSub(recv exchange.Exchange)
	BindPub(send exchange.Exchange) error

	LocalEvent(event func(data T) T)
	PrintInfo()

	GetId() uint64
	SetManager(*timestamp.TsManager[T, K, timestamp.DistributedClock[K], timestamp.LocalClock[K]])
}

// type timestamp.DistributedClock[K any] timestamp.DistributedClock[K]

type Message[T Payload, K ClockDataType] interface {
	String() string
	GetTs() timestamp.DistributedClock[K]
	GetData() T
	GetId() uint64
}

// func MessageToVersion[T Payload, K any, U timestamp.DistributedClock[K]](msg Message[T, K, U]) *timestamp.Version[T, K] {
// 	return timestamp.NewVersion[T, K](msg.GetData(), msg.GetTs(), msg.GetId())
// }
