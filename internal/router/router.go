package router

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"wbLvL0/internal/router/handlers"
	"wbLvL0/internal/storage"
)

const (
	homePageUrl = "/"
)

func InitRouter(st storage.Storage, log *slog.Logger) http.Handler {
	router := gin.Default()

	router.LoadHTMLGlob("html/*")
	router.Static("/static", "./static")

	h := handlers.New(st, log)
	router.GET(homePageUrl, h.GetOrder)

	return router
}
