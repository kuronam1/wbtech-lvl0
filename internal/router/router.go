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

	router.LoadHTMLFiles("render/html")
	router.Static("/static/css", "./static/css")

	h := handlers.New(st, log)
	router.GET(homePageUrl, h.GetOrder)

	return router
}
