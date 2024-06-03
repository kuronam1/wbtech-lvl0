package storage

import (
	"log/slog"
	"wbLvL0/internal/storage/orders"
	"wbLvL0/pkg/cache"
)

type Store struct {
	OrderRepo orders.Repository
	Cache     cache.Caсher
	Logger    *slog.Logger
}
