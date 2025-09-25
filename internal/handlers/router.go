package handlers

import (
	"net/http"

	"github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
)

type Router struct {
	Router *ginext.Engine
	store  *storage.Storage
}

func New(router *ginext.Engine, store *storage.Storage) *Router {
	return &Router{
		Router: router,
		store:  store,
	}
}

func (r *Router) Routers() {
	r.Router.POST("/add_notice", r.addNotice)
}

func (r *Router) addNotice(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{"error": "order UID is required"})

	notice := 0
	c.
		r.store.AddNotice()
}
