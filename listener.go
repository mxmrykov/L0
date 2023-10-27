package main

import (
	"fmt"
	"github.com/nats-io/stan.go"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type OrderModel struct {
	OrderUid string
}

const listenerID = "event-listener"

func cacheInput(message []byte) {
	//var currentOrder Order
	//err := json.Unmarshal(message, &currentOrder)
	//if err != nil {
	//	fmt.Println("cannot parse array")
	//}
	//fmt.Println(currentOrder)
	fmt.Println(string(message))
}

func main() {
	sc, err := stan.Connect("test-cluster", listenerID,
		stan.Pings(10, 3),
		stan.SetConnectionLostHandler(func(_ stan.Conn, err error) {
			log.Fatalf("Connection lost: %v", err)
		}))
	if err != nil {
		log.Fatalf("Error while connection: %v", err)
	}
	go func() {
		fmt.Println("Cluster 'test-cluster' connected")
		_, err := sc.Subscribe("main", func(m *stan.Msg) {
			fmt.Println("Received message:")
			cacheInput(m.Data)
		})
		if err != nil {
			log.Fatalf("Error while subscribing: %v", err)
		}
		fmt.Println("Chanel 'main' listener subscribed")
	}()
	fmt.Println("Connection closed")
	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel, syscall.SIGINT, syscall.SIGTERM)
	<-signalChanel
	err = sc.Close()
	if err != nil {
		log.Fatalf("Error while unsubscribing: %v", err)
	}
}
