package orders

import (
	"context"
	"wbLvL0/internal/models"
	"wbLvL0/internal/service/orders/db/dbModels"
)

type Repository interface {
	CreateFullOrder(order models.Order) error
	CreateOrder(ctx context.Context, order models.Order, deliveryID int, paymentID int) error
	GetOneOrder(uid string) (dbModels.Order, error)
	GetAllOrders() ([]dbModels.Order, error)
	CreatePayment(ctx context.Context, model models.Payment) (int, error)
	GetOnePayment(id int) (dbModels.Payment, error)
	CreateDelivery(ctx context.Context, model models.Delivery) (int, error)
	GetOneDelivery(id int) (dbModels.Delivery, error)
	CreateItem(ctx context.Context, item models.Item, orderUID string) error
	GetOrderItems(orderUID string) ([]dbModels.Item, error)
	CheckCache() (bool, error)
}
