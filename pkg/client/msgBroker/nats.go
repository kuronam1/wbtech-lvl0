package msgBroker

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"wbLvL0/internal/config"
	"wbLvL0/internal/storage/orders/models"
)

type Broker interface {
	Publish(subject string, data []byte) error
	Subscribe(ctx context.Context, subject string, cb func(data models.Order) error) error
}

func NewClient(cfg config.NatsStream) (stan.Conn, error) {
	conn, err := stan.Connect(cfg.ClusterID, cfg.ClientID, stan.NatsURL(nats.DefaultURL))
	if err != nil {
		return nil, fmt.Errorf("[nats-straming] cannot connect to nats: %w", err)
	}

	return conn, nil
}
