package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"wbLvL0/internal/models"
	"wbLvL0/internal/service/orders/db/dbModels"
	"wbLvL0/internal/service/orders/payment"
	"wbLvL0/pkg/client/postgreSQL"
)

type repository struct {
	Client postgreSQL.Client
}

func (r *repository) CreatePayment(ctx context.Context, model dbModels.Payment) (int, error) {
	var id int
	if err := r.Client.QueryRow(ctx, CreatePaymentQuery,
		model.PaymentID,
		model.Transaction,
		model.RequestId,
		model.Currency,
		model.Provider,
		model.Amount,
		model.PaymentDt,
		model.Bank,
		model.DeliveryCost,
		model.GoodsTotal,
		model.CustomFee,
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

func (r *repository) GetOnePayment(ctx context.Context, id int) (dbModels.Payment, error) {
	dbPayment := models.DBPayment{PaymentID: id}
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
		return models.DBPayment{}, err
	}

	return dbPayment, nil
}

func (r *repository) GetAllPayments(ctx context.Context) ([]dbModels.Payment, error) {
	rows, err := r.Client.Query(ctx, Se)
	if err != nil {
		return nil, err
	}

	payments := make([]models.DBPayment, 0)
	for rows.Next() {
		var OnePayment models.DBPayment

		if err := rows.Scan(
			&OnePayment.PaymentID,
			&OnePayment.Transaction,
			&OnePayment.RequestId,
			&OnePayment.Currency,
			&OnePayment.Provider,
			&OnePayment.Amount,
			&OnePayment.PaymentDt,
			&OnePayment.Bank,
			&OnePayment.DeliveryCost,
			&OnePayment.GoodsTotal,
			&OnePayment.CustomFee,
		); err != nil {
			return nil, err
		}

		payments = append(payments, OnePayment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return fullOrders, nil
}

func NewPayRepository(client postgreSQL.Client) payment.Repository {
	return &repository{
		Client: client,
	}
}
