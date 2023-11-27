package app

import (
	"fmt"
	"github.com/mxmrykov/L0/config"
	"github.com/mxmrykov/L0/internal/caches"
	nats "github.com/mxmrykov/L0/internal/nats"
	"github.com/mxmrykov/L0/internal/orders/controller"
	"github.com/mxmrykov/L0/internal/orders/generate"
	"github.com/mxmrykov/L0/internal/repository"
	"github.com/mxmrykov/L0/pkg/http"
	"github.com/mxmrykov/L0/pkg/postgres"
	"log"
	"time"
)

func Run(cfg *config.Config) {

	ns := nats.NewNats(&cfg.Nats)

	fmt.Println("Nats server started, connection registered")

	postgresConnection, err := postgres.Connect(&cfg.PG)

	if err != nil {
		fmt.Printf("Error at Postgre connection: %v", err)
	}
	defer postgresConnection.Close()

	repo := repository.NewRepository(postgresConnection)
	dbCreationError := repo.CreateTable()

	if dbCreationError != nil {
		fmt.Printf("Error at table creating: %v", dbCreationError)
	}

	orderCache := caches.NewCache(repo)
	orderCache.Preload()

	go func() {
		for {
			order := generate.GenerateOrder()
			fmt.Println("Order sent")
			err := ns.Publish(*order)

			if err != nil {
				fmt.Printf("Error at publishing: %v\n", err)
			}
			time.Sleep(30 * time.Second)
		}
	}()

	go func() {
		for {
			mes, err := ns.Subscribe()
			fmt.Println("Order received")
			if err != nil {
				fmt.Printf("Error at subscribing: %v", err)
			}

			if err != nil {
				fmt.Printf("Error at Unmarshaling: %v", err)
			}

			orderCache.CreateCache(*mes)
			time.Sleep(30 * time.Second)
		}
	}()

	httpServer := http.NewServer()
	orderController := controller.NewOrderController(orderCache)

	serverStartingError := httpServer.Start(orderController.GetOrderController, orderController.GetAllOrders)

	if serverStartingError != nil {
		log.Fatalf("Error at server starting: %v", serverStartingError)
	}
}
