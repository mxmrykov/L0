package models

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
	SmId              int      `json:"sm_id"`
	DateCreated       string   `json:"date_created"`
	OofShard          string   `json:"oof_shard"`
}
