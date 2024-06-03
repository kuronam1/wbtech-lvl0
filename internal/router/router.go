package router

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"wbLvL0/internal/router/handlers"
	"wbLvL0/internal/storage"
)

const (
	homePageUrl  = "/"
	orderPageUrl = "/order"
)

func InitRouter(st storage.Storage, log *slog.Logger) http.Handler {
	router := gin.Default()

	h := handlers.New(st, log)
	router.GET(homePageUrl, h.HomePage)
	router.GET(orderPageUrl, h.GetOrder)

	return router
}
