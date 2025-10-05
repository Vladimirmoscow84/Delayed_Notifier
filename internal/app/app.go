package app

import (
	"log"

	"github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/cache"
	"github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/handlers"
	datasaver "github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/service/data_saver"
	"github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/storage"
	wbconfig "github.com/wb-go/wbf/config"
	wbgin "github.com/wb-go/wbf/ginext"
)

func Run() {
	//ctx := context.Background()
	cfg := wbconfig.New()
	err := cfg.Load("../.env")
	if err != nil {
		log.Fatalf("load cfg dissable %v", err)
	}
	databaseUri := cfg.GetString("DATABASE_URI")
	addr := cfg.GetString("SERVER_ADDRESS")

	wbRouter := wbgin.New()

	store, err := storage.New(databaseUri, rdAddr)
	if err != nil {
		log.Fatalf("dissable connet to storage %v", err)
	}
	cache := cache.New(redisUri)

	dataSaverService := datasaver.New(store, cache)

	router := handlers.New(wbRouter, store)
	router.Routers()

	err = router.Router.Run(addr)
	if err != nil {
		log.Fatalf("connet to server dissadled %v", err)
	}

}
