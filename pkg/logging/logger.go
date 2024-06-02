package logging

import (
	"log/slog"
	"os"
	"wbLvL0/internal/config"
)

const (
	EnvLocal = "local"
	EnvProd  = "prod"
)

func InitLogger(cfg *config.Config) *slog.Logger {
	var logger *slog.Logger

	switch cfg.Env {
	case EnvLocal:
		logger = slog.New(slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			}))
	case EnvProd:
		logger = slog.New(slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelInfo,
			}))
	}

	logger.Info("logger initialized", "env", cfg.Env)

	return logger
}
