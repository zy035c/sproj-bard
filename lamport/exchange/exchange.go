package exchange

import (
	"time"
)

type MessageBroker[M any] struct {
	channel          chan M
	delay_max        time.Duration
	delay_min        time.Duration
	packet_loss_perc float32
}

func (mb *MessageBroker[M]) Put(m M) {
	mb.channel <- m
}

func (mb *MessageBroker[M]) Get() M {
	return <-mb.channel
}
