package handlers

import (
	"net/http"
	"time"

	"github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/model"
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
	r.Router.POST("/notify", r.addNotice)
	//r.Router.Get("/notify/:id", r.getStatus)
	//r.Router.Delete("/notify/:id", r.deleteNotice)
}

func (r *Router) addNotice(c *gin.Context) {
	ctx := c.Request.Context()

	var notice model.Notice

	if err := c.ShouldBindJSON(&notice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invlid input data."})
		return
	}

	notice.DateCreated = time.Now()
	notice.SendAttempts = 5
	notice.SendStatus = "sheduled"

	id, err := r.store.AddNotice(ctx, notice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}
