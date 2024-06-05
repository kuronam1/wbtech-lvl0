package orders

import (
	"context"
	"github.com/jackc/pgx/v5"
	"wbLvL0/internal/storage/orders/db/dbModels"
	"wbLvL0/internal/storage/orders/models"
)

type Repository interface {
	GetFullOrderByUID(uid string) (models.Order, error)
	ParseOrderData(ord dbModels.Order, delivery dbModels.Delivery,
		payment dbModels.Payment, items []dbModels.Item) (models.Order, error)
	CreateFullOrder(order models.Order) error
	CheckCache() (bool, error)
	GetAllOrdersUIDs() ([]string, error)
	GetOrderItems(orderUID string) ([]dbModels.Item, error)
	CreateOrder(ctx context.Context, tx pgx.Tx, order models.Order, deliveryID int, paymentID int) error
	CreatePayment(ctx context.Context, tx pgx.Tx, model models.Payment) (int, error)
	CreateDelivery(ctx context.Context, tx pgx.Tx, model models.Delivery) (int, error)
	CreateItem(ctx context.Context, tx pgx.Tx, item models.Item, orderUID string) error
}
