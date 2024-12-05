package simulation

import (
	"lamport/utils"
	"sync"
	"time"
)

type SimConfig struct {
	ReadWriteRatio float32
	AvgInterval    time.Duration
	AvgDelay       time.Duration
	PLR            float32
}

func (conf *SimConfig) PoissonInterval() time.Duration {
	return time.Duration(utils.Poisson(utils.Reciprocal(float64(conf.AvgInterval))))
}

type EventType uint8

const (
	READ EventType = 0
	WRITE
)

type Event struct {
	Vid   int
	Epoch int
	Etype EventType
	Op    func()
	Mid   int
}

type EventNode struct {
	Idx   int
	Event *Event
	Next  *EventNode
	Prev  *EventNode
}

type ClientHistory struct {
	first *EventNode
	last  *EventNode
	mutex sync.Mutex
}

func (c *ClientHistory) Add(mid int, event *Event) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.first == nil {
		c.first = &EventNode{
			Event: event,
			Idx:   mid,
		}
		c.last = c.first
		return
	}

	runner := c.first
	for runner.Next != nil {
		runner = runner.Next
	}

	runner.Next = &EventNode{
		Event: event,
		Idx:   mid,
		Prev:  runner,
	}

	c.last = runner.Next
}

func (c *ClientHistory) GetLast(mid int) *Event {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	runner := c.last

	for runner != nil {
		if runner.Idx == mid {
			return runner.Event
		}
		runner = runner.Prev
	}
	return nil
}
