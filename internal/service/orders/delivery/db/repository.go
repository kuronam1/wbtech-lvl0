package db

import (
	"wbLvL0/internal/service/orders/delivery"
	"wbLvL0/pkg/client/postgreSQL"
)

type repository struct {
	Client postgreSQL.Client
}

func (r *repository) CreateDelivery() {
	//TODO implement me
	panic("implement me")
}

func (r *repository) GetOne() {
	//TODO implement me
	panic("implement me")
}

func (r *repository) GetAllDeliveries() {
	//TODO implement me
	panic("implement me")
}

func NewDeliveryRepository(client postgreSQL.Client) delivery.Repository {
	return &repository{
		Client: client,
	}
}
