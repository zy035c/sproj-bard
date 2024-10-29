package exchange

type Exchange interface {
	Put(t []byte) bool
	Get() []byte
	C() <-chan []byte
	GetId() uint64
}
