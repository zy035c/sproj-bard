package utils

type Promise struct {
	Channel chan any
}

func (promise Promise) Then(f func(param any) any) Promise {
	ok := <-promise.Channel
	res := f(ok)
	nextp := Promise{}
	nextp.Channel <- res
	return nextp
}

func (promise Promise) ChanTo(ch chan any) {
	defer close(promise.Channel)
	elem := <-promise.Channel
	ch <- elem
}
