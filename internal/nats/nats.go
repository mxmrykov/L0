package nats

import (
	"encoding/json"
	"fmt"
	"github.com/mxmrykov/L0/config"
	"github.com/mxmrykov/L0/internal/models"
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

	sc, err := stan.Connect(cfg.Cluster, cfg.Client,
		stan.Pings(10, 3),
		stan.SetConnectionLostHandler(func(_ stan.Conn, err error) {
			fmt.Printf("Connection lost: %v", err)
		}))
	if err != nil {
		fmt.Printf("Cannot connect to stan: %v", err)
		return nil
	}

	return &Nats{cfg, sc, nc}
}

func (ns *Nats) Publish(message *models.Order) error {

	ord, err := json.MarshalIndent(message, "", "\t")

	if err != nil {
		fmt.Printf("Error at marshaling new order: %v", err)
	}

	return ns.sc.Publish(ns.config.Topic, ord)
}

func (ns *Nats) Subscribe() (*models.Order, error) {

	ch := make(chan *models.Order)

	_, err := ns.sc.Subscribe(ns.config.Topic, func(mes *stan.Msg) {

		var msg *models.Order

		err := json.Unmarshal(mes.Data, &msg)

		if err != nil {
			fmt.Printf("Error at Unmarshaling: %v", err)
		}

		ch <- msg
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
