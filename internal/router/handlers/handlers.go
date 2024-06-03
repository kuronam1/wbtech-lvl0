package handlers

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"wbLvL0/internal/storage"
)

type Handler interface {
	HomePage(c *gin.Context)
	GetOrder(c *gin.Context)
}

func New(store storage.Storage, logger *slog.Logger) Handler {
	return &handler{
		store:  store,
		logger: logger,
	}
}

func (h *handler) RegisterRoutes(router *gin.Engine) {
	router.GET(homePageUrl)
	router.GET(orderPageUrl)
}

func (h *handler) HomePage(c *gin.Context) {
	if c.Request.URL.Path != homePageUrl {
		c.AbortWithStatusJSON(404, gin.H{
			"error": "page not found",
		})
	}
	c.HTML(http.StatusOK, "", gin.H{
		"With order": false,
	})
}

func (h *handler) GetOrder(c *gin.Context) {
	uid, found := c.GetQuery("uid")
	if !found {
		c.Redirect(http.StatusMovedPermanently, homePageUrl)
		return
	}

	order, err := h.store.GetOrderByUID(uid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "error while getting an order",
		})
		return
	}

	c.HTML(http.StatusOK, "", gin.H{
		"With order": true,
		"Order":      order,
	})
}
