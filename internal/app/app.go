package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"wbLvL0/internal/appErrors"
	"wbLvL0/internal/broker"
	"wbLvL0/internal/config"
	"wbLvL0/internal/router"
	"wbLvL0/internal/storage"
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

	BRClient, err := msgBroker.NewClient(cfg.NatsStream, logger)
	if err != nil {
		log.Fatal(fmt.Sprintf("[BrClient] error while initializing new nats connetction [%s]", err))
	}

	br := broker.NewBroker(BRClient, logger)
	logger.Info("[Run] msgBroker initialized")

	st := storage.New(PGClient, logger)
	st.MustGetCache()
	logger.Info("[Run] storage initialized")

	rt := router.InitRouter(st, logger)
	logger.Info("[Run] routes and html files initialized")

	ctxStan, cancelSubPub := context.WithCancel(context.Background())
	br.Subscribe(ctxStan, cfg.NatsStream.ClusterID, st.CreateOrder)
	logger.Info("[Run] listening to cluster: %s", cfg.NatsStream.ClusterID)

	/*	br.Publish(ctxStan, cfg.NatsStream.ClusterID)
		logger.Info("[Run] publishing to cluster: %s", cfg.NatsStream.ClusterID)*/

	httpServer := server.New(rt, cfg.HttpServer)
	logger.Info("[Run] server started on addr: %s:%s", cfg.HttpServer.Host, cfg.HttpServer.Port)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	select {
	case stop := <-interrupt:
		logger.Error(fmt.Sprintf("[Run] os signal: %s", stop.String()))
	case svErr := <-httpServer.Notify:
		logger.Error(fmt.Sprintf("[Run] http signal: %s", appErrors.WrapLogErr(svErr)))
	case brErr := <-br.Notify:
		logger.Error(fmt.Sprintf("[Run] msgBroker signal: %s", appErrors.WrapLogErr(brErr)))
	}
	cancelSubPub()
	logger.Info("[Run] closing sub and pub")
	PGClient.Close()
	logger.Info("[Run] closing PG client connection")
	_ = br.Conn.Close()
	logger.Info("[Run] closing msgBroker client connection")
	logger.Info("[Run] shutting down server")
	err = httpServer.Shutdown()
	if err != nil {
		logger.Error(fmt.Sprintf("[Server] Stopped - http.Server.Shutdown: %s", appErrors.WrapLogErr(err)))
	}
}
