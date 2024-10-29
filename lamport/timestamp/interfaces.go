package timestamp

type TsOrder uint8

const (
	AFT TsOrder = iota
	BEF
	CON
)

type DistributedClock[T any] interface {
	DefaultTsCmp(DistributedClock[T]) TsOrder
	Value() T
	Increment()
	Set(data T)
	Clone() DistributedClock[T]
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}

type DistClockJson[T any] struct {
	Data      T      `json:"data"`
	ClockType string `json:"type"`
}

func GetDistType[T any](type_ string) DistributedClock[T] {
	res := (interface{})(nil)
	if type_ == "vector" {
		res = &VectorClock{}
	} else if type_ == "lamport" {
		res = &LamportClock{}
	}
	return res.(DistributedClock[T])
}

type LocalClock[T any] interface {
	Forward()
	Adjust(DistributedClock[T]) error
	Snapshot() T
	SnapshotTS() DistributedClock[T]
}
