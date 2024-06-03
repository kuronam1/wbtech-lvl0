package server

import (
	"context"
	"net/http"
	"time"
	"wbLvL0/internal/config"
)

type App struct {
	Server       *http.Server
	Notify       chan error
	shutdownTime time.Duration
}

func New(handler http.Handler, cfg config.HttpServer) *App {

	s := &http.Server{
		Addr:         cfg.Host + ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	}

	app := &App{
		Server:       s,
		Notify:       make(chan error),
		shutdownTime: cfg.ShutdownTimeout,
	}

	app.Start()

	return app
}

func (app *App) Start() {
	go func() {
		app.Notify <- app.Server.ListenAndServe()
		close(app.Notify)
	}()
}

func (app *App) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), app.shutdownTime)
	defer cancel()

	return app.Server.Shutdown(ctx)
}
