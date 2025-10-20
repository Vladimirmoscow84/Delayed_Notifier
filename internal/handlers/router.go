package handlers

import (
	"context"

	"github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/model"
	rabbitmq "github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/rabbitMq"
	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
)

type dataSaver interface {
	SaveData(ctx context.Context, notice model.Notice) error
}
type statusGetter interface {
	GetStatusNotice(ctx context.Context, id string) (string, error)
}
type dataDeleter interface {
	DeleteData(ctx context.Context, id string) error
}
type Router struct {
	Router       *ginext.Engine
	dataSaver    dataSaver
	statusGetter statusGetter
	dataDeleter  dataDeleter
	rabbit       *rabbitmq.Client
}

func New(router *ginext.Engine, dataSaver dataSaver, statusGetter statusGetter, dataDeleter dataDeleter, rabbitClient *rabbitmq.Client) *Router {
	return &Router{
		Router:       router,
		dataSaver:    dataSaver,
		statusGetter: statusGetter,
		dataDeleter:  dataDeleter,
		rabbit:       rabbitClient,
	}
}

func (r *Router) Routers() {
	r.Router.POST("/notify", r.addNotice)          //создание уведомлений с датой и временем отправки
	r.Router.GET("/notify/:id", r.getStatus)       //получение статуса уведомления по  ID
	r.Router.DELETE("/notify/:id", r.deleteNotice) //отмена запланированного уведомления по ID
	r.Router.GET("/", func(c *gin.Context) { c.File("./web/index.html") })
	r.Router.Static("/static", "./web")
}
