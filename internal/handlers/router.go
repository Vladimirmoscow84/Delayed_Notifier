package handlers

import (
	"context"

	"github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/model"
	"github.com/wb-go/wbf/ginext"
)

type service interface {
	SaveData(ctx context.Context, notice model.Notice) (int, error)
}
type Router struct {
	Router  *ginext.Engine
	service service
}

func New(router *ginext.Engine, service service) *Router {
	return &Router{
		Router:  router,
		service: service,
	}
}

func (r *Router) Routers() {
	r.Router.POST("/notify", r.addNotice)          //создание уведомлений с датой и временем отправки
	r.Router.GET("/notify/:id", r.getStatus)       //получение статуса уведомления по  ID
	r.Router.DELETE("/notify/:id", r.deleteNotice) //отмена запланированного уведомления по ID
}
