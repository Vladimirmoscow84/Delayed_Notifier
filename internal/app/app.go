package app

import (
	"log"

	"github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/handlers"
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
	rdAddr := cfg.GetString("REDIS_ADDRESS")

	router := wbgin.New()

	store, err := storage.New(databaseUri, rdAddr)
	if err != nil {
		log.Fatalf("dissable connet to storage %v", err)
	}

	handls := handlers.New(router, store)
	handls.Routers()

	err = handls.Router.Run(addr)
	if err != nil {
		log.Fatalf("connet to server dissadled %v", err)
	}

}
