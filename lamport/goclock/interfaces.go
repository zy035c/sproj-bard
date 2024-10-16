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

type Machine interface {
	Start()
	Stop()
	Broadcast()
	Listen()
	LocalEvent()
}

type Message[T Payload, K any] struct {
	timestamp.Version[T, K, timestamp.DistributedClock[K]]
}
