package handlers

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"wbLvL0/internal/storage"
)

const (
	homePageUrl = "/"
)

type Handler interface {
	GetOrder(c *gin.Context)
}

func New(store storage.Storage, logger *slog.Logger) Handler {
	return &handler{
		store:  store,
		logger: logger,
	}
}

func (h *handler) GetOrder(c *gin.Context) {
	if c.Request.URL.Path != homePageUrl {
		c.AbortWithStatusJSON(404, gin.H{
			"error": "page not found",
		})
	}

	uid, found := c.GetQuery("uid")
	if !found {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"WithOrder": false,
		})
		return
	}

	order, err := h.store.GetOrderByUID(uid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "error while getting an order",
		})
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"WithOrder": true,
		"Order":     order,
	})
}
