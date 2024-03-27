package mypq

import (
	"fmt"
	"sync/atomic"
	"time"
)

type Item struct {
	Value    func() // value
	Priority uint64
	index    int
}

func (item *Item) GetPriority() uint64 {
	return item.Priority
}

func (item *Item) SetPriority(priority uint64) {
	item.Priority = priority
}

func (item *Item) Execute() {
	item.Value()
}

// PriorityQueue
type PriorityQueue struct {
	items []Item
	ch    chan struct{}
}

// heap.Interface
func (pq *PriorityQueue) Len() int {
	return len(pq.items)
}

func (pq *PriorityQueue) Less(i, j int) bool {
	return pq.items[i].Priority < pq.items[j].Priority
}

func (pq *PriorityQueue) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
	pq.items[i].index = i
	pq.items[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Item)
	item.index = len(pq.items)
	pq.items = append(pq.items, *item)
	pq.ch <- struct{}{}
}

func (pq *PriorityQueue) Pop() interface{} {
	old := pq.items
	n := len(old)
	item := old[n-1]
	item.index = -1
	pq.items = old[0 : n-1]
	return item
}

// constructor
func NewPriorityQueue() *PriorityQueue {
	return &PriorityQueue{ch: make(chan struct{}, 100)} // buffer size 100
}

func (pq *PriorityQueue) LoopAndPoll(ts *uint64) {
	timer := time.NewTimer(0)
	for {
		timer.Reset(time.Millisecond * 100)

		select {
		case <-pq.ch:
			cur_timestamp := atomic.LoadUint64(ts)
			item := pq.items[0]
			if item.GetPriority() == cur_timestamp || item.GetPriority() == cur_timestamp+1 {
				// convert the value to a function and execute it
				fmt.Println("[LoopAndPoll] Will execute item. Current timestamp", cur_timestamp, "; item's timestamp", item.GetPriority())
				item = pq.Pop().(Item)
				item.Execute()
			} else {
				pq.ch <- struct{}{}
			}
		case <-timer.C:
			// pass
		}
	}
}

// func main() {

// 	pq := PriorityQueue{}

// 	heap.Push(&pq, &Item{value: func() { println("foo") }, priority: 3})
// 	heap.Push(&pq, &Item{value: func() { println("bar") }, priority: 1})
// 	heap.Push(&pq, &Item{value: func() { println("bty") }, priority: 2})

// 	println("len", pq.Len())

// 	for pq.Len() > 0 {
// 		item := heap.Pop(&pq).(Item)
// 		item.Execute()
// 	}
// }
