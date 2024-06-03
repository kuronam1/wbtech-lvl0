package orders

import (
	"context"
	"github.com/jackc/pgx/v5"
	"wbLvL0/internal/storage/orders/db/dbModels"
	"wbLvL0/internal/storage/orders/models"
)

type Repository interface {
	CreateFullOrder(order models.Order) error
	CheckCache() (bool, error)
	GetOneOrder(uid string) (dbModels.Order, error)
	GetAllOrders() ([]dbModels.Order, error)
	GetOnePayment(id int) (dbModels.Payment, error)
	GetOneDelivery(id int) (dbModels.Delivery, error)
	GetOrderItems(orderUID string) ([]dbModels.Item, error)
	CreateOrder(ctx context.Context, tx pgx.Tx, order models.Order, deliveryID int, paymentID int) error
	CreatePayment(ctx context.Context, tx pgx.Tx, model models.Payment) (int, error)
	CreateDelivery(ctx context.Context, tx pgx.Tx, model models.Delivery) (int, error)
	CreateItem(ctx context.Context, tx pgx.Tx, item models.Item, orderUID string) error
}
