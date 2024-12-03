package simulation

import (
	"lamport/utils"
	"time"
)

type SimConfig struct {
	ReadWriteRatio float32
	AvgInterval    time.Duration
	AvgDelay       time.Duration
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
	Epoch int
	Etype EventType
	Op    func()
	Mid   int
}

type ClientHistory struct {
	IdxToEvent map[int]*Event
}
