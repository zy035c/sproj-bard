package test

type Intf[T uint16 | uint32] interface {
	Foo()
	Bar()
}

func TakeInnf[T uint16 | uint32](x Intf[T]) {

}

type valRec struct {
}

func (v valRec) Foo() {

}

func (v valRec) Bar() {

}

type ptrRec struct {
}

func (v *ptrRec) Foo() {

}

func (v *ptrRec) Bar() {

}

func MyF() {

	a := valRec{}
	TakeInnf[uint16](a)

	b := ptrRec{}
	TakeInnf[uint32](&b)

}
