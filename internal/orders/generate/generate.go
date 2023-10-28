package generate

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/mxmrykov/L0/internal/models"
	"math"
	"math/rand"
	"strconv"
	"time"
)

func GenerateOrder() *models.Order {
	orderId := hash32bit()
	var orderCount = 1 + rand.Intn(2)
	items := make([]models.Item, orderCount)
	for i := 0; i < orderCount; i += 1 {
		items[i] = models.Item{0, "", 0, "", "", 0, "", 0, 0, "", 0}
	}
	order := models.Order{
		orderId[:len(orderId)-15],
		"WBILMTESTTRACK",
		"WBIL",
		generateOrderDelivery(),
		models.Payment{"", "", "", "", 0, 0, "", 0, 0, 0},
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
	generateOrderItems(&order)
	generateOrderPayment(&order)
	return &order
	//ord, err := json.MarshalIndent(order, "", "\t")
	//if err != nil {
	//	fmt.Printf("Error while converting order to json: %v", err)
	//	return nil
	//}
	//return ord
}

func generateOrderPayment(order *models.Order) {
	currency := []string{"USD", "RUB", "EUR"}
	banks := []string{"sber", "alpha", "tinkoff"}
	var amount float32 = 0
	for i := range order.Items {
		amount += order.Items[i].TotalPrice
	}
	deliveryCost := float32(rand.Intn(1500))
	order.Payment.Transaction = order.OrderUid + order.CustomerId
	order.Payment.RequestId = ""
	order.Payment.Currency = currency[rand.Intn(len(currency))]
	order.Payment.Provider = "wbpay"
	order.Payment.Amount = amount + deliveryCost
	order.Payment.PaymentDt = uint32(1000000000 + rand.Intn(1000000000))
	order.Payment.Bank = banks[rand.Intn(len(banks))]
	order.Payment.DeliveryCost = uint32(deliveryCost)
	order.Payment.GoodsTotal = amount
	order.Payment.CustomFee = 0
}

func generateOrderDelivery() models.Delivery {
	names := []string{"Rykov Maxim", "Ivanov Ivan", "Random Random", "Maximov Maxim"}
	addresses := []string{"Ploshad Mira 15", "Pokryshkina 8", "Orshanskaya 3", "Red Square"}
	newDelivery := models.Delivery{
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

func generateOrderItems(order *models.Order) {
	for i := 0; i < len(order.Items); i += 1 {
		amount := float32(1500 + rand.Intn(10000))
		sale := float32(rand.Intn(50))
		total := ((100 - sale) / 100.0) * amount
		totalPrice := math.Round(float64(total)*10) / 10
		order.Items[i] = models.Item{
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

func hash32bit() string {
	sum := md5.Sum([]byte(strconv.Itoa(rand.Intn(150000))))
	return hex.EncodeToString(sum[:])
}
