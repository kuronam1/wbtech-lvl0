package broker

import (
	"log/slog"
	"wbLvL0/pkg/client/msgBroker"
)

type Stan struct {
	Conn   msgBroker.Client
	Logger *slog.Logger
	Notify chan error
}
