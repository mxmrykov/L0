package orders

//func main() {
//	subs, err := stan.Connect("test-cluster", publisherID,
//		stan.Pings(10, 3),
//		stan.SetConnectionLostHandler(func(_ stan.Conn, err error) {
//			log.Fatalf("Connection lost: %v", err)
//		}))
//	if err != nil {
//		log.Fatalf("Error while connection: %v", err)
//	}
//	go func() {
//
//		fmt.Println("Cluster 'test-cluster' connected")
//		for {
//			err = subs.Publish("main", generateOrder())
//			if err != nil {
//				log.Fatalf("Error while publishing: %v", err)
//			}
//			fmt.Println("Message sent")
//			time.Sleep(60 * time.Second)
//		}
//	}()
//	signalChanel := make(chan os.Signal, 1)
//	signal.Notify(signalChanel, syscall.SIGINT, syscall.SIGTERM)
//	<-signalChanel
//	err = subs.Close()
//	if err != nil {
//		log.Fatalf("Error while unsubscribing: %v", err)
//	}
//}
