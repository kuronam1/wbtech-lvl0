package db

import (
	"wbLvL0/internal/service/orders/items"
	"wbLvL0/pkg/client/postgreSQL"
)

type repository struct {
	Client postgreSQL.Client
}

func (r *repository) CreateItem() {
	//TODO implement me
	panic("implement me")
}

func (r *repository) GetOne() {
	//TODO implement me
	panic("implement me")
}

func (r *repository) GetAllItemsByOrderID() {
	//TODO implement me
	panic("implement me")
}

func NewItemRepository(client postgreSQL.Client) items.Repository {
	return &repository{
		Client: client,
	}
}
