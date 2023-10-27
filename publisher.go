package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
	"log"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const publisherID = "event-publisher"

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string  `json:"transaction"`
	RequestId    string  `json:"request_id"`
	Currency     string  `json:"currency"`
	Provider     string  `json:"provider"`
	Amount       float32 `json:"amount"`
	PaymentDt    uint32  `json:"payment_dt"`
	Bank         string  `json:"bank"`
	DeliveryCost uint32  `json:"delivery_cost"`
	GoodsTotal   float32 `json:"goods_total"`
	CustomFee    uint32  `json:"custom_fee"`
}

type Item struct {
	ChrtId      uint32  `json:"chrt_id"`
	TrackNumber string  `json:"track_number"`
	Price       uint16  `json:"price"`
	Rid         string  `json:"rid"`
	Name        string  `json:"name"`
	Sale        uint16  `json:"sale"`
	Size        string  `json:"size"`
	TotalPrice  float32 `json:"total_price"`
	NmId        uint32  `json:"nm_id"`
	Brand       string  `json:"brand"`
	Status      uint16  `json:"status"`
}

type Order struct {
	OrderUid          string   `json:"order_uid"`
	TrackNumber       string   `json:"track_number"`
	Entry             string   `json:"entry"`
	Delivery          Delivery `json:"delivery"`
	Payment           Payment  `json:"payment"`
	Items             []Item   `json:"items"`
	Locale            string   `json:"locale"`
	InternalSignature string   `json:"internal_signature"`
	CustomerId        string   `json:"customer_id"`
	DeliveryService   string   `json:"delivery_service"`
	ShardKey          string   `json:"shard_key"`
	SmId              uint32   `json:"sm_id"`
	DateCreated       string   `json:"date_created"`
	OofShard          string   `json:"oof_shard"`
}

func (order *Order) generateOrderPayment(pm *Payment) {
	currency := []string{"USD", "RUB", "EUR"}
	banks := []string{"sber", "alpha", "tinkoff"}
	var amount float32 = 0
	for i := range order.Items {
		amount += order.Items[i].TotalPrice
	}
	deliveryCost := float32(rand.Intn(1500))
	pm.Transaction = order.OrderUid + order.CustomerId
	pm.RequestId = ""
	pm.Currency = currency[rand.Intn(len(currency))]
	pm.Provider = "wbpay"
	pm.Amount = amount + deliveryCost
	pm.PaymentDt = uint32(1000000000 + rand.Intn(1000000000))
	pm.Bank = banks[rand.Intn(len(banks))]
	pm.DeliveryCost = uint32(deliveryCost)
	pm.GoodsTotal = amount
	pm.CustomFee = 0
}

func hash32bit() string {
	sum := md5.Sum([]byte(strconv.Itoa(rand.Intn(150000))))
	return hex.EncodeToString(sum[:])
}

func generateOrderDelivery() Delivery {
	names := []string{"Rykov Maxim", "Ivanov Ivan", "Random Random", "Maximov Maxim"}
	addresses := []string{"Ploshad Mira 15", "Pokryshkina 8", "Orshanskaya 3", "Red Square"}
	newDelivery := Delivery{
		names[rand.Intn(len(names))],
		"+" + strconv.Itoa(1000000000+rand.Intn(8000000000)),
		strconv.Itoa(100000 + rand.Intn(150000)),
		"Moscow",
		addresses[rand.Intn(len(addresses))],
		"Moscow",
		"test@gmail.com",
	}
	return newDelivery
}

func (order *Order) generateOrderItems(orderItems *[]Item) {
	for i := 0; i < len(*orderItems); i += 1 {
		amount := float32(1500 + rand.Intn(10000))
		sale := float32(rand.Intn(50))
		total := ((100 - sale) / 100.0) * amount
		totalPrice := math.Round(float64(total)*10) / 10
		(*orderItems)[i] = Item{
			uint32(rand.Intn(1000000)),
			order.TrackNumber,
			uint16(amount),
			hash32bit()[:len(hash32bit())-15] + order.CustomerId,
			"Mascaras",
			uint16(sale),
			"0",
			float32(totalPrice),
			uint32(rand.Intn(1000000)),
			"Vivienne Sabo",
			202,
		}
	}
}

func generateOrder() []byte {
	orderId := hash32bit()
	var orderCount = 1 + rand.Intn(2)
	items := make([]Item, orderCount)
	for i := 0; i < orderCount; i += 1 {
		items[i] = Item{0, "", 0, "", "", 0, "", 0, 0, "", 0}
	}
	order := &Order{
		orderId[:len(orderId)-15],
		"WBILMTESTTRACK",
		"WBIL",
		generateOrderDelivery(),
		Payment{"", "", "", "", 0, 0, "", 0, 0, 0},
		items,
		"en",
		"",
		"text",
		"meest",
		"9",
		99,
		time.Now().Format(time.RFC3339),
		"1",
	}
	order.generateOrderItems(&order.Items)
	order.generateOrderPayment(&order.Payment)
	fmt.Println("NEW ORDER GENERATED")
	ord, err := json.MarshalIndent(order, "", "\t")
	if err != nil {
		fmt.Printf("Error while converting order to json: %v", err)
		return nil
	}
	return ord
}

func main() {
	subs, err := stan.Connect("test-cluster", publisherID,
		stan.Pings(10, 3),
		stan.SetConnectionLostHandler(func(_ stan.Conn, err error) {
			log.Fatalf("Connection lost: %v", err)
		}))
	if err != nil {
		log.Fatalf("Error while connection: %v", err)
	}
	go func() {

		fmt.Println("Cluster 'test-cluster' connected")
		for {
			err = subs.Publish("main", generateOrder())
			if err != nil {
				log.Fatalf("Error while publishing: %v", err)
			}
			fmt.Println("Message sent")
			time.Sleep(60 * time.Second)
		}
	}()
	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel, syscall.SIGINT, syscall.SIGTERM)
	<-signalChanel
	err = subs.Close()
	if err != nil {
		log.Fatalf("Error while unsubscribing: %v", err)
	}
}
