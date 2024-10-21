package goclock

import "lamport/timestamp"

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
	Broadcast(Message[T, K])
	Listen()
	LocalEvent(event func(data T) T)
	SetSend(send []chan Message[T, K]) error
	PrintInfo()
}

type ClockAbbr[K any] timestamp.DistributedClock[K]

type Message[T Payload, K any] interface {
	String() string
	GetTs() timestamp.DistributedClock[K]
	GetData() T
	GetId() uint64
}

// func MessageToVersion[T Payload, K any, U ClockAbbr[K]](msg Message[T, K, U]) *timestamp.Version[T, K, U] {
// 	return timestamp.NewVersion[T, K, U](msg.GetData(), msg.GetTs(), msg.GetId())
// }
