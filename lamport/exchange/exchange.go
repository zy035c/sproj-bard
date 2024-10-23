package exchange

import (
	"fmt"
	"lamport/utils"
	"time"
)

type ExchangeImpl[M any] struct {
	channel          chan M
	delay_max        time.Duration
	delay_min        time.Duration
	packet_loss_perc float32
	id               uint64
}

func NewExchangeImpl[M any](delay_max, delay_min time.Duration, packet_loss_perc float32, broker_size uint32, id uint64) (*ExchangeImpl[M], error) {
	if packet_loss_perc > 1 || packet_loss_perc < 0 {
		return nil, fmt.Errorf("! packet_loss_perc must be a 0~1 float")
	}
	if broker_size == 0 {
		broker_size = 512
	}
	return &ExchangeImpl[M]{
		channel:          make(chan M, broker_size),
		delay_max:        delay_max,
		delay_min:        delay_min,
		packet_loss_perc: packet_loss_perc,
		id:               id,
	}, nil
}

func (mb *ExchangeImpl[M]) Put(m M) bool {
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

func (mb *ExchangeImpl[M]) Get() M {
	return <-mb.channel
}

func (mb *ExchangeImpl[M]) GetId() uint64 {
	return mb.id
}

func (mb *ExchangeImpl[M]) C() <-chan M {
	return mb.channel
}

type ConsistentHash[M any] struct {
	ExchangeImpl[M]
	// TODO
}
