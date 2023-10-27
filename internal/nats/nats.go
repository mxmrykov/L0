package nats

import (
	"fmt"
	"github.com/mxmrykov/L0/config"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"time"
)

type Nats struct {
	config *config.Nats
	sc     stan.Conn
	nc     *nats.Conn
}

func NewNats(cfg *config.Nats) *Nats {
	natsUrl := fmt.Sprintf("nats://%s:%s", cfg.Host, cfg.Port)

	nc, err := nats.Connect(natsUrl)
	if err != nil {
		fmt.Printf("Cannot connect to nats: %v", err)
		return nil
	}

	sc, err := stan.Connect(cfg.Cluster, cfg.Client)
	if err != nil {
		fmt.Printf("Cannot connect to stan: %v", err)
		return nil
	}

	return &Nats{cfg, sc, nc}
}

func (ns *Nats) Publish(message []byte) error {
	return ns.sc.Publish(ns.config.Topic, message)
}

func (ns *Nats) Subscribe() (*[]byte, error) {

	ch := make(chan *[]byte)

	_, err := ns.sc.Subscribe(ns.config.Topic, func(mes *stan.Msg) {
		ch <- &mes.Data
	})

	if err != nil {
		fmt.Printf("Error at subscription: %v", err)
	}

	select {
	case message := <-ch:
		return message, nil
	case <-time.After(60 * time.Second):
		return nil, stan.ErrTimeout
	}

}
