package app

import (
	"context"
	"log"
	"time"

	"github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/cache"
	"github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/handlers"
	"github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/model"
	"github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/notifier/email"
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
	eHost := cfg.GetString("EMAIL_HOST")
	ePort := cfg.GetString("EMAIL_PORT")
	eUser := cfg.GetString("EMAIL_USER")
	ePass := cfg.GetString("EMAIL_PASS")
	eFrom := cfg.GetString("EMAIL_FROM")
	eTo := cfg.GetString("EMAIL_TO")

	wbRouter := wbgin.New("release")

	//--------ТЕСТ ПОЧТЫ ---------
	client := email.New(eHost, ePort, eUser, ePass, eFrom, eTo)
	sendTime, _ := time.Parse(time.RFC3339, "2025-10-15T07:03:00Z")
	notice := model.Notice{
		Id:           123,
		Body:         "sosi",
		DateCreated:  sendTime,
		SendDate:     sendTime,
		SendAttempts: 3,
		SendStatus:   "fdhbvdhfd",
	}

	client.Send(notice)
	if err != nil {
		log.Fatalf("Ошибка отправки письма: %v", err)
	}

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
