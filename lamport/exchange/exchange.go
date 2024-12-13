package exchange

import (
	"fmt"
	"lamport/utils"
	"time"
)

type SingleBroker struct {
	Channel          chan []byte
	DelayMax         time.Duration
	DelayMin         time.Duration
	PacketLossPerctg float32
	Id               uint64
}

type PoissonBroker struct {
	SingleBroker
	DelayAvg time.Duration
}

func NewSingleBroker(delay_max, delay_min time.Duration, packet_loss_perc float32, broker_size uint32, id uint64) (*SingleBroker, error) {
	if packet_loss_perc > 1 || packet_loss_perc < 0 {
		return nil, fmt.Errorf("! packet_loss_perc must be a 0~1 float")
	}
	if broker_size == 0 {
		broker_size = 512
	}
	return &SingleBroker{
		Channel:          make(chan []byte, broker_size),
		DelayMax:         delay_max,
		DelayMin:         delay_min,
		PacketLossPerctg: packet_loss_perc,
		Id:               id,
	}, nil
}

func NewPoissonBroker(avg_delay time.Duration, avg_packet_loss_perc float32, broker_size uint32, id uint64) (*PoissonBroker, error) {

	if broker_size == 0 {
		broker_size = 512
	}
	return &PoissonBroker{
		SingleBroker: SingleBroker{
			Channel:          make(chan []byte, broker_size),
			DelayMax:         0,
			DelayMin:         0,
			PacketLossPerctg: avg_packet_loss_perc,
			Id:               id,
		},
		DelayAvg: avg_delay,
	}, nil
}

func (mb *SingleBroker) Put(m []byte) bool {
	if mb.PacketLossPerctg != 0 {
		if utils.RandomFloat32(0, 1) < mb.PacketLossPerctg {
			return false
		}
	}

	randt := utils.RandomInt64(int64(mb.DelayMin), int64(mb.DelayMax))
	time.Sleep(time.Duration(randt))
	mb.Channel <- m
	return true
}

func (mb *SingleBroker) Get() []byte {
	return <-mb.Channel
}

func (mb *SingleBroker) GetId() uint64 {
	return mb.Id
}

func (mb *SingleBroker) C() <-chan []byte {
	return mb.Channel
}

type ConsistentHash struct {
	SingleBroker
	// TODO
}

func (mb *PoissonBroker) Put(m []byte) bool {
	if mb.PacketLossPerctg > 0 {
		if utils.RandomFloat32(0, 1) < mb.PacketLossPerctg {
			return false
		}
	}
	if mb.DelayAvg > 0 {
		time.Sleep(time.Duration(utils.Poisson(utils.Reciprocal(float64(mb.DelayAvg)))))
	}
	mb.Channel <- m
	return true
}
