package cache

import (
	"log/slog"
	"sync"
)

type Cache struct {
	sync.RWMutex
	Items  map[string]Item
	Logger *slog.Logger
}

type Item interface{}
