package service

import (
	"context"
	"github.com/jinzhu/copier"
	"log"
	"log/slog"
	"wbLvL0/internal/errors"
	"wbLvL0/internal/models"
	"wbLvL0/internal/service/orders"
	"wbLvL0/internal/service/orders/db"
	"wbLvL0/pkg/cache"
	"wbLvL0/pkg/client/postgreSQL"
)

type Service struct {
	OrderRepo orders.Repository
	Cache     cache.Caсher
	Logger    *slog.Logger
}

func New(PGClient postgreSQL.Client, logger *slog.Logger) *Service {

	return &Service{
		OrderRepo: db.NewRepository(PGClient),
		Cache:     cache.InitCache(),
		Logger:    logger,
	}
}

func (s *Service) MustGetCache(ctx context.Context) {
	cacheFounded, err := s.OrderRepo.CheckCache(ctx)
	if err != nil {
		s.Logger.Error("error while getting orders", errors.WrapLogErr(err))
		log.Fatal("cannot get cache")
	}

	if cacheFounded {
		AllOrders, err := s.OrderRepo.GetAllOrders(ctx)
		if err != nil {
			s.Logger.Error("error while getting orders", errors.WrapLogErr(err))
			log.Fatal("cannot get cache")
		}

		for _, order := range AllOrders {
			var realOrder models.Order
			if err := copier.Copy(&realOrder, &order); err != nil {
				s.Logger.Error("error while copy order", errors.WrapLogErr(err))
				log.Fatal("cannot get cache")
			}

			oneDelivery, err := s.OrderRepo.GetOneDelivery(ctx, order.OrderID)
			if err != nil {
				s.Logger.Error("error while getting order's delivery", errors.WrapLogErr(err))
				log.Fatal("cannot get cache")
			}

			var realDelivery models.Delivery
			if err := copier.Copy(&realDelivery, &oneDelivery); err != nil {
				s.Logger.Error("error while copy order's delivery", errors.WrapLogErr(err))
				log.Fatal("cannot get cache")
			}

			onePayment, err := s.OrderRepo.GetOnePayment(ctx, order.OrderID)
			if err != nil {
				s.Logger.Error("error while getting order's payment", errors.WrapLogErr(err))
				log.Fatal("cannot get cache")
			}

			var realPayment models.Payment
			if err := copier.Copy(&realPayment, &onePayment); err != nil {
				s.Logger.Error("error while copy order's payment", errors.WrapLogErr(err))
				log.Fatal("cannot get cache")
			}

			orderItems, err := s.OrderRepo.GetOrderItems(ctx, order.OrderUid)
			if err != nil {
				s.Logger.Error("error while getting order's items", errors.WrapLogErr(err))
				log.Fatal("cannot get cache")
			}

			var realItems []models.Item
			if err := copier.Copy(&realItems, &orderItems); err != nil {
				s.Logger.Error("error while copy order's items", errors.WrapLogErr(err))
				log.Fatal("cannot get cache")
			}

			realOrder.Delivery = realDelivery
			realOrder.Payment = realPayment
			realOrder.Items = realItems

			s.Cache.Set(order.OrderUid, order)
		}
	}
}
