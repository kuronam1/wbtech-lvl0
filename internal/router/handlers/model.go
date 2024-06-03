package handlers

import (
	"log/slog"
	"wbLvL0/internal/storage"
)

type handler struct {
	store  storage.Storage
	logger *slog.Logger
}
