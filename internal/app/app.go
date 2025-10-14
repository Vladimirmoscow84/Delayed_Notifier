package app

import (
	"log"

	"github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/cache"
	"github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/handlers"
	rabbitmq "github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/rabbitMq"
	datadeleter "github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/service/data_deleter"
	datasaver "github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/service/data_saver"
	statusgetter "github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/service/status_getter"
	"github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/storage"
	wbconfig "github.com/wb-go/wbf/config"
	wbgin "github.com/wb-go/wbf/ginext"
)

func Run() {
	//ctx := context.Background()
	cfg := wbconfig.New()
	err := cfg.Load("", "../.env", "")
	if err != nil {
		log.Fatalf("load cfg dissable %v", err)
	}
	databaseUri := cfg.GetString("DATABASE_URI")
	addr := cfg.GetString("SERVER_ADDRESS")
	redisUri := cfg.GetString("REDIS_URI")
	rabbitUri := cfg.GetString("RABBIT_URI")

	wbRouter := wbgin.New("release")

	store, err := storage.New(databaseUri)
	if err != nil {
		log.Fatalf("dissable connet to storage %v", err)
	}
	cache := cache.New(redisUri)

	rabbitCfg := rabbitmq.Config{
		RabbitUri:    rabbitUri,
		Exchange:     "main_exchange",
		Exchangetype: "direct",
		Queue:        "main_queue",
		RoutingKey:   "notices",
		DLX:          "dlx_exchange",
		DLQ:          "dlq_queue",
	}

	dataSaverService := datasaver.New(store, cache)
	statusGetterService := statusgetter.New(cache)
	dataDeleterService := datadeleter.New(cache)
	rabbitClient, err := rabbitmq.New(rabbitCfg)
	if err != nil {
		log.Fatalf("failed connect to RabbitMQ: %v", err)
	}
	defer rabbitClient.Close()

	router := handlers.New(wbRouter, dataSaverService, statusGetterService, dataDeleterService, rabbitClient)
	router.Routers()

	err = router.Router.Run(addr)
	if err != nil {
		log.Fatalf("connet to server dissadled %v", err)
	}

}
