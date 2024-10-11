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
}

type LocalClock[T any] interface {
	Forward()
	Adjust(DistributedClock[T]) error
	Snapshot() T
}
