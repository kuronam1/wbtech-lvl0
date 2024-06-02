package dbModels

import "time"

type Order struct {
	OrderID           int `copier:"-"`
	OrderUid          string
	TrackNumber       string
	Entry             string
	DeliveryID        int
	PaymentID         int
	Locale            string
	InternalSignature string
	CustomerId        string
	DeliveryService   string
	Shardkey          string
	SmId              int
	DateCreated       time.Time
	OofShard          string
}

type Delivery struct {
	DeliveryID int `copier:"-"`
	Name       string
	Phone      string
	Zip        string
	City       string
	Address    string
	Region     string
	Email      string
}

type Payment struct {
	PaymentID    int `copier:"-"`
	Transaction  string
	RequestId    string
	Currency     string
	Provider     string
	Amount       int
	PaymentDt    int
	Bank         string
	DeliveryCost int
	GoodsTotal   int
	CustomFee    int
}

type Item struct {
	ItemID      int    `copier:"-"`
	OrderUID    string `copier:"-"`
	ChrtId      int
	TrackNumber string
	Price       int
	Rid         string
	Name        string
	Sale        int
	Size        string
	TotalPrice  int
	NmId        int
	Brand       string
	Status      int
}
