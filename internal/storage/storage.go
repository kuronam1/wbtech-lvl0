package storage

import (
	"fmt"
	"github.com/jinzhu/copier"
	"log"
	"log/slog"
	"wbLvL0/internal/errors"
	"wbLvL0/internal/storage/orders/db"
	"wbLvL0/internal/storage/orders/models"
	"wbLvL0/pkg/cache"
	"wbLvL0/pkg/client/postgreSQL"
)

type Storage interface {
	MustGetCache()
	GetOrderByUID(uid string) (models.Order, error)
}

func New(PGClient postgreSQL.Client, logger *slog.Logger) Storage {

	return &Store{
		OrderRepo: db.NewRepository(PGClient, logger),
		Cache:     cache.InitCache(),
		Logger:    logger,
	}
}

func (s *Store) MustGetCache() {
	cacheFounded, err := s.OrderRepo.CheckCache()
	if err != nil {
		s.Logger.Error("error while getting orders", errors.WrapLogErr(err))
		log.Fatal("cannot get cache")
	}

	if cacheFounded {
		AllOrders, err := s.OrderRepo.GetAllOrders()
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

			oneDelivery, err := s.OrderRepo.GetOneDelivery(order.OrderID)
			if err != nil {
				s.Logger.Error("error while getting order's delivery", errors.WrapLogErr(err))
				log.Fatal("cannot get cache")
			}

			var realDelivery models.Delivery
			if err := copier.Copy(&realDelivery, &oneDelivery); err != nil {
				s.Logger.Error("error while copy order's delivery", errors.WrapLogErr(err))
				log.Fatal("cannot get cache")
			}

			onePayment, err := s.OrderRepo.GetOnePayment(order.OrderID)
			if err != nil {
				s.Logger.Error("error while getting order's payment", errors.WrapLogErr(err))
				log.Fatal("cannot get cache")
			}

			var realPayment models.Payment
			if err := copier.Copy(&realPayment, &onePayment); err != nil {
				s.Logger.Error("error while copy order's payment", errors.WrapLogErr(err))
				log.Fatal("cannot get cache")
			}

			orderItems, err := s.OrderRepo.GetOrderItems(order.OrderUid)
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

func (s *Store) GetOrderByUID(uid string) (models.Order, error) {
	order, _ := s.Cache.Get(uid)

	o, ok := order.(models.Order)
	if !ok {
		s.Logger.Error("error while casting", "order", order)
		return models.Order{}, fmt.Errorf("error while casting")
	}

	return o, nil
}
