package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jinzhu/copier"
	"log/slog"
	"time"
	"wbLvL0/internal/appErrors"
	"wbLvL0/internal/storage/orders"
	"wbLvL0/internal/storage/orders/db/dbModels"
	"wbLvL0/internal/storage/orders/models"
	"wbLvL0/pkg/client/postgreSQL"
)

const (
	RWTimeout       = 5 * time.Second
	GetCacheTimeout = 5 * time.Second
)

type repository struct {
	Client postgreSQL.Client
	Logger *slog.Logger
}

func NewRepository(client postgreSQL.Client, logger *slog.Logger) orders.Repository {
	return &repository{
		Client: client,
		Logger: logger,
	}
}

func (r *repository) CreateFullOrder(order models.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), RWTimeout)
	defer cancel()

	tx, err := r.Client.Begin(ctx)
	if err != nil {
		r.Logger.Error("[repository] error while creating transaction: ", appErrors.WrapLogErr(err))
		return err
	}
	defer tx.Rollback(ctx)

	paymentID, err := r.CreatePayment(ctx, tx, order.Payment)
	if err != nil {
		r.Logger.Error("[repository] error while creating payment: ", appErrors.WrapLogErr(err))
		return err
	}

	deliveryID, err := r.CreateDelivery(ctx, tx, order.Delivery)
	if err != nil {
		r.Logger.Error("[repository] error while creating delivery: ", appErrors.WrapLogErr(err))
		return err
	}

	for _, item := range order.Items {
		if err := r.CreateItem(ctx, tx, item, order.OrderUid); err != nil {
			r.Logger.Error("[repository] error while creating Item: ", appErrors.WrapLogErr(err))
			return err
		}
	}

	if err := r.CreateOrder(ctx, tx, order, deliveryID, paymentID); err != nil {
		r.Logger.Error("[repository] error while creating dbOrder: ", appErrors.WrapLogErr(err))
		return err
	}

	return tx.Commit(ctx)
}

func (r *repository) CreateOrder(ctx context.Context, tx pgx.Tx, order models.Order, deliveryID int, paymentID int) error {
	var (
		dbOrder dbModels.Order
		id      int
	)
	if err := copier.Copy(&dbOrder, &order); err != nil {
		return err
	}

	if err := tx.QueryRow(ctx, CreateOrderQuery,
		dbOrder.OrderUid,
		dbOrder.TrackNumber,
		dbOrder.Entry,
		deliveryID,
		paymentID,
		dbOrder.Locale,
		dbOrder.InternalSignature,
		dbOrder.CustomerId,
		dbOrder.DeliveryService,
		dbOrder.Shardkey,
		dbOrder.SmId,
		dbOrder.DateCreated,
		dbOrder.OofShard,
	).Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			newErr := fmt.Errorf(
				fmt.Sprintf(
					"SQL Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where,
				))
			return newErr

		}
		return err
	}

	return nil
}

func (r *repository) GetOneOrder(uid string) (dbModels.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), RWTimeout)
	defer cancel()

	order := dbModels.Order{OrderUid: uid}
	if err := r.Client.QueryRow(ctx, SelectOrderQuery, uid).Scan(
		&order.OrderID,
		&order.TrackNumber,
		&order.Entry,
		&order.DeliveryID,
		&order.PaymentID,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerId,
		&order.DeliveryService,
		&order.Shardkey,
		&order.SmId,
		&order.DateCreated,
		&order.OofShard,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			newErr := fmt.Errorf(
				fmt.Sprintf(
					"SQL Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where,
				))
			return dbModels.Order{}, newErr

		}
		return dbModels.Order{}, err
	}

	return order, nil
}

func (r *repository) GetAllOrders() ([]dbModels.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), RWTimeout)
	defer cancel()

	rows, err := r.Client.Query(ctx, SelectOrdersQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fullOrders := make([]dbModels.Order, 0)
	for rows.Next() {
		var order dbModels.Order

		if err := rows.Scan(
			order.OrderID,
			order.OrderUid,
			order.TrackNumber,
			order.Entry,
			order.DeliveryID,
			order.PaymentID,
			order.Locale,
			order.InternalSignature,
			order.CustomerId,
			order.DeliveryService,
			order.Shardkey,
			order.SmId,
			order.DateCreated,
			order.OofShard,
		); err != nil {
			return nil, err
		}

		fullOrders = append(fullOrders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return fullOrders, nil
}

func (r *repository) CreatePayment(ctx context.Context, tx pgx.Tx, model models.Payment) (int, error) {
	var (
		dbPayment dbModels.Payment
		id        int
	)
	if err := copier.Copy(&dbPayment, &model); err != nil {
		return 0, err
	}

	if err := tx.QueryRow(ctx, CreatePaymentQuery,
		dbPayment.Transaction,
		dbPayment.RequestId,
		dbPayment.Currency,
		dbPayment.Provider,
		dbPayment.Amount,
		dbPayment.PaymentDt,
		dbPayment.Bank,
		dbPayment.DeliveryCost,
		dbPayment.GoodsTotal,
		dbPayment.CustomFee,
	).Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			newErr := fmt.Errorf(
				fmt.Sprintf(
					"SQL Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where,
				))
			return 0, newErr

		}
		return 0, err
	}

	return id, nil
}

func (r *repository) GetOnePayment(id int) (dbModels.Payment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), RWTimeout)
	defer cancel()

	dbPayment := dbModels.Payment{PaymentID: id}
	if err := r.Client.QueryRow(ctx, SelectPaymentQuery, id).Scan(
		&dbPayment.Transaction,
		&dbPayment.RequestId,
		&dbPayment.Currency,
		&dbPayment.Provider,
		&dbPayment.Amount,
		&dbPayment.PaymentDt,
		&dbPayment.Bank,
		&dbPayment.DeliveryCost,
		&dbPayment.GoodsTotal,
		&dbPayment.CustomFee,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			newErr := fmt.Errorf(
				fmt.Sprintf(
					"SQL Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where,
				))
			return dbModels.Payment{}, newErr

		}
		return dbModels.Payment{}, err
	}

	return dbPayment, nil
}

func (r *repository) CreateDelivery(ctx context.Context, tx pgx.Tx, model models.Delivery) (int, error) {
	var (
		dbDelivery dbModels.Delivery
		id         int
	)
	if err := copier.Copy(&dbDelivery, &model); err != nil {
		return 0, err
	}

	if err := tx.QueryRow(ctx, CreateDeliveryQuery,
		dbDelivery.Name,
		dbDelivery.Phone,
		dbDelivery.Zip,
		dbDelivery.City,
		dbDelivery.Address,
		dbDelivery.Region,
		dbDelivery.Email,
	).Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			newErr := fmt.Errorf(
				fmt.Sprintf(
					"SQL Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where,
				))
			return 0, newErr

		}
		return 0, err
	}

	return id, nil
}

func (r *repository) GetOneDelivery(id int) (dbModels.Delivery, error) {
	ctx, cancel := context.WithTimeout(context.Background(), RWTimeout)
	defer cancel()

	delivery := dbModels.Delivery{DeliveryID: id}
	if err := r.Client.QueryRow(ctx, SelectDeliveryQuery, id).Scan(
		&delivery.Name,
		&delivery.Phone,
		&delivery.Zip,
		&delivery.City,
		&delivery.Address,
		&delivery.Region,
		&delivery.Email,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			newErr := fmt.Errorf(
				fmt.Sprintf(
					"SQL Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where,
				))
			return dbModels.Delivery{}, newErr

		}
		return dbModels.Delivery{}, err
	}

	return delivery, nil
}

func (r *repository) CreateItem(ctx context.Context, tx pgx.Tx, item models.Item, orderUID string) error {
	var dbItem dbModels.Item
	if err := copier.Copy(&dbItem, &item); err != nil {
		return err
	}

	r.Logger.Debug(fmt.Sprintf("%v", dbItem))

	_, err := tx.Exec(ctx, CreateItemQuery,
		orderUID,
		item.ChrtId,
		item.TrackNumber,
		item.Price,
		item.Rid,
		item.Name,
		item.Sale,
		item.Size,
		item.TotalPrice,
		item.NmId,
		item.Brand,
		item.Status,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			newErr := fmt.Errorf(
				fmt.Sprintf(
					"SQL Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where,
				))
			return newErr

		}
		return err
	}

	return nil
}

func (r *repository) GetOrderItems(orderUID string) ([]dbModels.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), RWTimeout)
	defer cancel()

	rows, err := r.Client.Query(ctx, SelectItemsQuery, orderUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]dbModels.Item, 0)
	for rows.Next() {
		var item dbModels.Item
		if err := rows.Scan(
			&item.ItemID,
			&item.OrderUID,
			&item.ChrtId,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmId,
			&item.Brand,
			&item.Status,
		); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				newErr := fmt.Errorf(
					fmt.Sprintf(
						"SQL Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where,
					))
				return nil, newErr

			}
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *repository) CheckCache() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), GetCacheTimeout)
	defer cancel()

	var flag bool
	if err := r.Client.QueryRow(ctx, CheckCacheQuery).Scan(&flag); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			newErr := fmt.Errorf(
				fmt.Sprintf(
					"SQL Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where,
				))
			return flag, newErr

		}
		return flag, err
	}
	return flag, nil
}
