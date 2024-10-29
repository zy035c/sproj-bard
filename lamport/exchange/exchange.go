package exchange

import (
	"fmt"
	"lamport/utils"
	"time"
)

type SingleBroker struct {
	channel          chan []byte
	delay_max        time.Duration
	delay_min        time.Duration
	packet_loss_perc float32
	id               uint64
}

func NewSingleBroker(delay_max, delay_min time.Duration, packet_loss_perc float32, broker_size uint32, id uint64) (*SingleBroker, error) {
	if packet_loss_perc > 1 || packet_loss_perc < 0 {
		return nil, fmt.Errorf("! packet_loss_perc must be a 0~1 float")
	}
	if broker_size == 0 {
		broker_size = 512
	}
	return &SingleBroker{
		channel:          make(chan []byte, broker_size),
		delay_max:        delay_max,
		delay_min:        delay_min,
		packet_loss_perc: packet_loss_perc,
		id:               id,
	}, nil
}

func (mb *SingleBroker) Put(m []byte) bool {
	if mb.packet_loss_perc != 0 {
		if utils.RandomFloat32(0, 1) < mb.packet_loss_perc {
			return false
		}
	}

	randt := utils.RandomInt64(int64(mb.delay_min), int64(mb.delay_max))
	time.Sleep(time.Duration(randt))
	mb.channel <- m
	return true
}

func (mb *SingleBroker) Get() []byte {
	return <-mb.channel
}

func (mb *SingleBroker) GetId() uint64 {
	return mb.id
}

func (mb *SingleBroker) C() <-chan []byte {
	return mb.channel
}

type ConsistentHash struct {
	SingleBroker
	// TODO
}
