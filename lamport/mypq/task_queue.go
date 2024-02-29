package mypq

import "sync"

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
	lock  sync.Mutex
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
	pq.lock.Lock()
	defer pq.lock.Unlock()

	item := x.(*Item)
	item.index = len(pq.items)
	pq.items = append(pq.items, *item)
}

func (pq *PriorityQueue) Pop() interface{} {
	pq.lock.Lock()
	defer pq.lock.Unlock()

	old := pq.items
	n := len(old)
	item := old[n-1]
	item.index = -1
	pq.items = old[0 : n-1]
	return item
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
