package app

import (
	"context"
	"fmt"
	"log"
	"time"
	"wbLvL0/internal/broker"
	"wbLvL0/internal/config"
	"wbLvL0/internal/service"
	"wbLvL0/internal/service/orders/db"
	"wbLvL0/pkg/client/msgBroker"
	"wbLvL0/pkg/client/postgreSQL"
	"wbLvL0/pkg/logging"
)

const (
	connectPGTimeout = 5 * time.Second
)

func Run(cfg *config.Config) {
	logger := logging.InitLogger(cfg)

	ctxPG, cancelPG := context.WithTimeout(context.Background(), connectPGTimeout)
	defer cancelPG()

	PGClient, err := postgreSQL.NewClient(ctxPG, cfg.PG)
	if err != nil {
		log.Fatal(fmt.Sprintf("[PGClient] error while initializing new PGClient [%s]", err))
	}

	repo := db.NewRepository(PGClient)

	BrClient, err := msgBroker.NewClient(cfg.NatsStream)
	if err != nil {
		log.Fatal(fmt.Sprintf("[BrClient] error while initializing new BrClient [%s]", err))
	}

	stan := broker.NewBroker(BrClient, logger)

	ctxStun, cancel := context.WithCancel(context.Background())
	go func() {
		if err := stan.Subscribe(); err != nil {

		}
	}()

	Service := service.New(PGClient, InMemCache, BrClient, logger)

	ctxCache, cancel := context.WithTimeout(context.Background(), GetCacheTimeout)
	defer cancel()
	Service.MustGetCache(ctxCache)
}
