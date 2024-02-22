package main

import (
	"fmt"
	"sync"
	"time"
)

type Event struct {
	ID     int
	Clock  int
	Action string
}

type DistributedStorage struct {
	Data  map[string]string
	Clock int
	Mutex sync.Mutex
}

func (ds *DistributedStorage) SetValue(key, value string) {
	ds.Mutex.Lock()
	defer ds.Mutex.Unlock()

	ds.Clock++
	ds.Data[key] = value
	fmt.Printf("Set key '%s' to value '%s' at time %d\n", key, value, ds.Clock)
}

func (ds *DistributedStorage) GetValue(key string) (string, bool) {
	ds.Mutex.Lock()
	defer ds.Mutex.Unlock()

	value, exists := ds.Data[key]
	if !exists {
		return "", false
	}
	fmt.Printf("Get value '%s' for key '%s' at time %d\n", value, key, ds.Clock)
	return value, true
}

func (ds *DistributedStorage) LamportClock() int {
	ds.Mutex.Lock()
	defer ds.Mutex.Unlock()

	ds.Clock++
	return ds.Clock
}

func main() {
	storage := DistributedStorage{
		Data:  make(map[string]string),
		Clock: 0,
		Mutex: sync.Mutex{},
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		storage.SetValue("foo", "bar")
	}()

	go func() {
		time.Sleep(200 * time.Millisecond)
		storage.SetValue("baz", "qux")
	}()

	time.Sleep(300 * time.Millisecond)

	storage.GetValue("foo")
	storage.GetValue("baz")
}
