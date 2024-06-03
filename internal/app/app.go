package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"wbLvL0/internal/broker"
	"wbLvL0/internal/config"
	"wbLvL0/internal/errors"
	"wbLvL0/internal/router"
	"wbLvL0/internal/storage"
	"wbLvL0/internal/storage/orders/db"
	"wbLvL0/pkg/client/msgBroker"
	"wbLvL0/pkg/client/postgreSQL"
	"wbLvL0/pkg/logging"
	"wbLvL0/pkg/server"
)

func Run(cfg *config.Config) {
	logger := logging.InitLogger(cfg)
	logger.Info("[Logger] initialized")

	PGClient, err := postgreSQL.NewClient(cfg.PG)
	if err != nil {
		log.Fatal(fmt.Sprintf("[PGClient] error while initializing new DB connetion pool [%s]", err))
	}
	logger.Info("[PGClient] connection established")

	repo := db.NewRepository(PGClient, logger)
	logger.Info("[Repository] initialized")

	BrClient, err := msgBroker.NewClient(cfg.NatsStream)
	if err != nil {
		log.Fatal(fmt.Sprintf("[BrClient] error while initializing new nats connetction [%s]", err))
	}
	logger.Info("[BRClient] connection established")

	br := broker.New(BrClient, logger)
	logger.Info("[Broker] initialized")

	ctxStan, cancelStan := context.WithCancel(context.Background())
	//go st.Publish()
	go br.Subscribe(ctxStan, cfg.NatsStream.ClusterID, repo.CreateFullOrder)
	logger.Info("[Broker] listening to cluster: %s", cfg.NatsStream.ClusterID)

	st := storage.New(PGClient, logger)
	st.MustGetCache()
	logger.Info("[Storage] storage initialized")

	rt := router.InitRouter(st, logger)
	logger.Info("[Router] routes and html files initialized")

	httpServer := server.New(rt, cfg.HttpServer)
	logger.Info("[Server] server started on addr: %s:%s", cfg.HttpServer.Host, cfg.HttpServer.Port)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	select {
	case stop := <-interrupt:
		logger.Error(fmt.Sprintf("[Run] os signal: %s", stop.String()))
	case svErr := <-httpServer.Notify:
		logger.Error(fmt.Sprintf("[Run] http signal: %s", errors.WrapLogErr(svErr)))
	case brErr := <-br.Notify:
		logger.Error(fmt.Sprintf("[Run] broker signal: %s", errors.WrapLogErr(brErr)))
	}

	defer func() {
		cancelStan()
		_ = PGClient.Close
		_ = br.Conn.Close()
	}()
	logger.Error("[Run] shutting down server")
	err = httpServer.Shutdown()
	if err != nil {
		logger.Error(fmt.Sprintf("[Server] Stopped - http.Server.Shutdown: %s", errors.WrapLogErr(err)))
	}
}
