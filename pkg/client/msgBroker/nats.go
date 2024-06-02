package msgBroker

import (
	"fmt"
	"github.com/nats-io/stan.go"
	"wbLvL0/internal/config"
	"wbLvL0/internal/models"
)

type Broker interface {
	Publish(subject string, data []byte) error
	Subscribe(subject string, cb func(data models.Order)) (func() error, error)
}

func NewClient(cfg config.NatsStream) (stan.Conn, error) {
	conn, err := stan.Connect(cfg.ClusterID, cfg.ClientID, stan.NatsURL(cfg.NatsUrl))
	if err != nil {
		return nil, fmt.Errorf("[nats-straming] cannot connect to nats: %w", err)
	}

	return conn, nil
}
