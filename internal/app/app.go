package app

import (
	"log"

	"github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/cache"
	"github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/handlers"
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

	wbRouter := wbgin.New("release")

	store, err := storage.New(databaseUri)
	if err != nil {
		log.Fatalf("dissable connet to storage %v", err)
	}
	cache := cache.New(redisUri)

	dataSaverService := datasaver.New(store, cache)
	statusGetterService := statusgetter.New(cache)
	dataDeleterService := datadeleter.New(cache)

	router := handlers.New(wbRouter, dataSaverService, statusGetterService, dataDeleterService)
	router.Routers()

	err = router.Router.Run(addr)
	if err != nil {
		log.Fatalf("connet to server dissadled %v", err)
	}

}
