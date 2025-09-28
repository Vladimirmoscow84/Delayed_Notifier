package handlers

import (
	"github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/storage"
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
	r.Router.POST("/notify", r.addNotice)    //создание уведомлений с датой и временем отправки
	r.Router.GET("/notify/:id", r.getStatus) //получение статуса уведомления
	//r.Router.DELETE("/notify/:id", r.deleteNotice) //отмена запланированного уведомления
}
