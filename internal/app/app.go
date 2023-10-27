package app

import (
	"fmt"
	"github.com/mxmrykov/L0/config"
	"github.com/mxmrykov/L0/internal/caches"
	nats "github.com/mxmrykov/L0/internal/nats"
	"github.com/mxmrykov/L0/internal/orders/generate"
	"github.com/mxmrykov/L0/pkg/http"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(cfg *config.Config) {

	ns := nats.NewNats(&cfg.Nats)

	webServer := http.Server{}

	serverStartingError := webServer.Start(&config.HTTP{})

	if serverStartingError != nil {
		log.Fatalf("Error at server starting: %v", serverStartingError)
	}

	fmt.Println("Nats server started, connection registered")

	go func() {
		for {
			mes, err := ns.Subscribe()

			if err != nil {
				fmt.Printf("Error at subscribing: %v", err)
			}

			if err != nil {
				fmt.Printf("Error at Unmarshaling: %v", err)
			}

			fmt.Println(cfg.Topic, ": ", mes.OrderUid)

			time.Sleep(50 * time.Second)
		}
	}()

	go func() {
		for {
			order := generate.GenerateOrder()

			err := ns.Publish(order)

			if err != nil {
				fmt.Printf("Error at publishing: %v\n", err)
			}

			orderCache := caches.NewCache()
			orderCache.CreateCache(order)

			time.Sleep(50 * time.Second)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

}

func messageSender(ns *config.Nats) {

}
