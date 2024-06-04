package models

import "time"

type Order struct {
	OrderUid          string `json:"order_uid" validator:"required"`
	TrackNumber       string `json:"track_number" validator:"required"`
	Entry             string `json:"entry" validator:"required"`
	Delivery          `json:"delivery" copier:"-" validator:"required"`
	Payment           `json:"payment" copier:"-" validator:"required"`
	Items             []Item    `json:"items" copier:"-" validator:"required,ne=0"`
	Locale            string    `json:"locale" validator:"required"`
	InternalSignature string    `json:"internal_signature" validator:"required"`
	CustomerId        string    `json:"customer_id" validator:"required"`
	DeliveryService   string    `json:"delivery_service" validator:"required"`
	Shardkey          string    `json:"shardkey" validator:"required"`
	SmId              int       `json:"sm_id" validator:"required"`
	DateCreated       time.Time `json:"date_created" validator:"required"`
	OofShard          string    `json:"oof_shard" validator:"required"`
}

type Payment struct {
	Transaction  string `json:"transaction" validator:"required"`
	RequestId    string `json:"request_id" validator:"required"`
	Currency     string `json:"currency" validator:"required"`
	Provider     string `json:"provider" validator:"required"`
	Amount       int    `json:"amount" validator:"required"`
	PaymentDt    int    `json:"payment_dt" validator:"required"`
	Bank         string `json:"bank" validator:"required"`
	DeliveryCost int    `json:"delivery_cost" validator:"required"`
	GoodsTotal   int    `json:"goods_total" validator:"required"`
	CustomFee    int    `json:"custom_fee" validator:"required"`
}

type Delivery struct {
	Name    string `json:"name" validator:"required"`
	Phone   string `json:"phone" validator:"required"`
	Zip     string `json:"zip" validator:"required"`
	City    string `json:"city" validator:"required"`
	Address string `json:"address" validator:"required"`
	Region  string `json:"region" validator:"required"`
	Email   string `json:"email" validator:"required"`
}

type Item struct {
	ChrtId      int    `json:"chrt_id" validator:"required"`
	TrackNumber string `json:"track_number" validator:"required"`
	Price       int    `json:"price" validator:"required"`
	Rid         string `json:"rid" validator:"required"`
	Name        string `json:"name" validator:"required"`
	Sale        int    `json:"sale" validator:"required"`
	Size        string `json:"size" validator:"required"`
	TotalPrice  int    `json:"total_price" validator:"required"`
	NmId        int    `json:"nm_id" validator:"required"`
	Brand       string `json:"brand" validator:"required"`
	Status      int    `json:"status" validator:"required"`
}
