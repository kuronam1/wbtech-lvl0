package storage

import (
	"fmt"
	"log"
	"log/slog"
	"wbLvL0/internal/appErrors"
	"wbLvL0/internal/storage/orders/db"
	"wbLvL0/internal/storage/orders/models"
	"wbLvL0/pkg/cache"
	"wbLvL0/pkg/client/postgreSQL"
)

type Storage interface {
	MustGetCache()
	GetOrderByUID(uid string) (models.Order, error)
	CreateOrder(order models.Order) error
}

func New(PGClient postgreSQL.Client, logger *slog.Logger) Storage {

	store := &Store{
		OrderRepo: db.NewRepository(PGClient, logger),
		Cache:     cache.InitCache(),
		Logger:    logger,
	}

	store.MustGetCache()

	return store
}

func (s *Store) MustGetCache() {
	cacheFounded, err := s.OrderRepo.CheckCache()
	if err != nil {
		s.Logger.Error("error while getting orders", appErrors.WrapLogErr(err))
		log.Fatal("cannot get cache")
	}

	if cacheFounded {
		AllOrdersUIDs, err := s.OrderRepo.GetAllOrdersUIDs()
		if err != nil {
			s.Logger.Error("error while getting orders", appErrors.WrapLogErr(err))
			log.Fatal("cannot get cache")
		}

		for _, uid := range AllOrdersUIDs {
			order, err := s.OrderRepo.GetFullOrderByUID(uid)
			if err != nil {
				s.Logger.Error("[Store] getting orders", appErrors.WrapLogErr(err))
				log.Fatal("cannot get cache")
			}

			s.Cache.Set(order.OrderUid, order)
		}
	}
}

func (s *Store) GetOrderByUID(uid string) (models.Order, error) {
	order, found := s.Cache.Get(uid)
	switch found {
	case false:
		order, err := s.OrderRepo.GetFullOrderByUID(uid)
		if err != nil {
			s.Logger.Error("[Store] getting order from db:", appErrors.WrapLogErr(err))
			return models.Order{}, err
		}

		return order, nil
	default:
		o, ok := order.(models.Order)
		if !ok {
			s.Logger.Error("error while casting", "order", order)
			return models.Order{}, fmt.Errorf("error while casting")
		}

		return o, nil
	}
}

func (s *Store) CreateOrder(order models.Order) error {
	if err := s.OrderRepo.CreateFullOrder(order); err != nil {
		s.Logger.Error("[s.CreateOrder] error while creating order: ", appErrors.WrapLogErr(err))
		return err
	}

	s.Cache.Set(order.OrderUid, order)
	return nil
}
