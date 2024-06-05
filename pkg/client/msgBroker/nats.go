package msgBroker

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"log/slog"
	"wbLvL0/internal/config"
)

type Client interface {
	Publish(subject string, data []byte) error
	Subscribe(subject string, cb stan.MsgHandler, opts ...stan.SubscriptionOption) (stan.Subscription, error)
	Close() error
}

func NewClient(cfg config.NatsStream, logger *slog.Logger) (stan.Conn, error) {
	conn, err := stan.Connect(cfg.ClusterID, cfg.ClientID, stan.NatsURL(nats.DefaultURL))
	if err != nil {
		return nil, fmt.Errorf("[nats-straming] cannot connect to nats: %w", err)
	}

	return conn, nil
}
