package app

import (
	"context"
	"log"

	"github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/cache"
	"github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/handlers"
	rabbitmq "github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/rabbitMq"
	datadeleter "github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/service/data_deleter"
	datasaver "github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/service/data_saver"
	statusgetter "github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/service/status_getter"
	amqp "github.com/rabbitmq/amqp091-go"
	wbconfig "github.com/wb-go/wbf/config"
	wbgin "github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/redis"
)

func Run() {
	//ctx := context.Background()
	cfg := wbconfig.New()
	err := cfg.Load("../config.yaml", "../.env", "")
	if err != nil {
		log.Fatalf("load cfg dissable %v", err)
	}
	//databaseUri := cfg.GetString("DATABASE_URI")
	addr := cfg.GetString("SERVER_ADDRESS")
	redisUri := cfg.GetString("REDIS_URI")
	rabbitUri := cfg.GetString("RABBIT_URI")

	wbRouter := wbgin.New("release")

	// store, err := storage.New(databaseUri,)
	// if err != nil {
	// 	log.Fatalf("dissable connet to storage %v", err)
	// }
	rd := redis.New(redisUri, "", 0)
	store := cache.NewCache(rd)

	rabbitCfg := rabbitmq.Config{
		RabbitUri:    rabbitUri,
		Exchange:     "main_exchange",
		Exchangetype: "direct",
		Queue:        "main_queue",
		RoutingKey:   "notices",
		DLX:          "dlx_exchange",
		DLQ:          "dlq_queue",
	}

	dataSaverService := datasaver.New(store)
	statusGetterService := statusgetter.New(store)
	dataDeleterService := datadeleter.New(store)
	rabbitClient, err := rabbitmq.New(rabbitCfg)
	if err != nil {
		log.Fatalf("failed connect to RabbitMQ: %v", err)
	}
	defer rabbitClient.Close()

	router := handlers.New(wbRouter, dataSaverService, statusGetterService, dataDeleterService, rabbitClient)
	router.Routers()

	go func() {
		log.Printf("Server is running 1")
		err = router.Router.Run(addr)
		if err != nil {
			log.Fatalf("connet to server dissadled %v", err)
		}
		log.Printf("Server is running 2")
	}()

	handler := func(msg amqp.Delivery) {
		log.Printf("Получено сообщение: %s", string(msg.Body))
	}
	err = rabbitClient.ConsumeDLQWithWorkers(context.Background(), 5, handler)
	if err != nil {
		log.Fatalf("failed consume with RabbitMQ: %v", err)
	}
	select {}
}
