package utils

import (
	"sync"
	"time"
)

type ThreadPool struct {
	size int
	wg   sync.WaitGroup
	sema chan struct{}
	freq time.Duration
}

func (t *ThreadPool) Init(size int) {
	t.size = size
	t.freq = time.Millisecond * 400
	t.sema = make(chan struct{}, 1024)
	for i := 0; i < t.size; i++ {
		t.sema <- struct{}{}
	}
}

func (t *ThreadPool) Wg() *sync.WaitGroup {
	return &t.wg
}

func (t *ThreadPool) SetWg(n int) {
	t.wg.Add(n)
}

func (t *ThreadPool) Submit(f func() any) {
	go func() {
		for {
			select {
			case <-t.sema:
				f()
				t.sema <- struct{}{}
				return
			case <-time.After(t.freq):
			}
		}
	}()
}
