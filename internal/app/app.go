package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/cache"
	"github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/handlers"
	"github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/notifier/email"
	"github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/notifier/telegram"
	rabbitmq "github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/rabbitMq"
	datadeleter "github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/service/data_deleter"
	datasaver "github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/service/data_saver"
	statusgetter "github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/service/status_getter"
	amqp "github.com/rabbitmq/amqp091-go"
	wbconfig "github.com/wb-go/wbf/config"
	wbgin "github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/redis"
)

// структура полезной нагрузки из RabbitMQ.
type NotifyMessage struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func Run() {
	cfg := wbconfig.New()
	err := cfg.Load("../config.yaml", "../.env", "")
	if err != nil {
		log.Fatalf("[main]load cfg dissable %v", err)
	}

	addr := cfg.GetString("SERVER_ADDRESS")
	redisUri := cfg.GetString("REDIS_URI")
	rabbitUri := cfg.GetString("RABBIT_URI")
	eHost := cfg.GetString("EMAIL_HOST")
	ePort := cfg.GetString("EMAIL_PORT")
	eUser := cfg.GetString("EMAIL_USER")
	ePass := cfg.GetString("EMAIL_PASS")
	eFrom := cfg.GetString("EMAIL_FROM")
	eTo := cfg.GetString("EMAIL_TO")
	tToken := cfg.GetString("TELEGRAM_TOKEN")
	tChatID := cfg.GetInt("TELEGRAM_CHAT_ID")
	// Инициализация уведомлений для email.
	emailClient := email.New(eHost, ePort, eUser, ePass, eFrom, eTo)

	// Инициализация уведомлений для telegram.
	var telegramClient *telegram.Client
	if tToken != "" && tChatID != 0 {
		tclient, err := telegram.New(tToken, int64(tChatID))
		if err != nil {
			log.Printf("[telega] init failed: %v", err)
		} else {
			telegramClient = tclient
		}
	} else {
		log.Printf("[telega] config missing, telegram disabled")
	}

	// Инициализация кэш и redis.
	rd := redis.New(redisUri, "", 0)
	store := cache.NewCache(rd)

	// Конфиг RabbitMQ.
	rabbitCfg := rabbitmq.Config{
		RabbitUri:    rabbitUri,
		Exchange:     "main_exchange",
		Exchangetype: "direct",
		Queue:        "main_queue",
		RoutingKey:   "notices",
		DLX:          "dlx_exchange",
		DLQ:          "dlq_queue",
	}

	//Инициализация сервисов приложения.
	dataSaverService := datasaver.New(store)
	statusGetterService := statusgetter.New(store)
	dataDeleterService := datadeleter.New(store)
	rabbitClient, err := rabbitmq.New(rabbitCfg)
	if err != nil {
		log.Fatalf("failed connect to RabbitMQ: %v", err)
	}
	defer rabbitClient.Close()

	wbRouter := wbgin.New("release")
	router := handlers.New(wbRouter, dataSaverService, statusGetterService, dataDeleterService, rabbitClient)
	router.Routers()

	//Запуск HTTP сервера.
	go func() {
		log.Printf("Server is running ")
		err = router.Router.Run(addr)
		if err != nil {
			log.Fatalf("[main]connet to server dissadled %v", err)
		}
	}()

	//Контекст завершения.
	ctx, cancel := context.WithCancel(context.Background())
	go handleShutdown(cancel)

	//Обработчик сообщений из DLQ.
	handler := func(msg amqp.Delivery) {
		var data NotifyMessage
		if err := json.Unmarshal(msg.Body, &data); err != nil {
			log.Printf("[handler] bad message JSON: %v", err)
			_ = rabbitClient.Nack(msg)
			return
		}

		log.Printf("[handler] received message: %+v", data)
		fmt.Println(data.Subject, "-----", data.Body)
		//Отправка уведомлений в телегу и на почту.
		if emailClient != nil {
			if err := emailClient.Send(data.Subject, data.Body); err != nil {
				log.Printf("[email] send error: %v", err)
			}
		}

		if telegramClient != nil {
			if err := telegramClient.Send(data.Body); err != nil {
				log.Printf("[telega] send error: %v", err)
			}
		}

		//Подтверждение удачной обработки.
		_ = rabbitClient.Ack(msg)
	}

	//Запуск воркеров для DLQ.
	err = rabbitClient.ConsumeDLQWithWorkers(ctx, 5, handler)
	if err != nil {
		log.Fatalf("failed consume from RabbitMQ DLQ: %v", err)
	}

	log.Println("[main] Application started. Waiting for messages...")

	<-ctx.Done()
	log.Println("[main] graceful shutdown complete.")
}

// Graceful shutdown при SIGINT/SIGTERM.
func handleShutdown(cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("[main] shutdown signal received")
	cancel()
}
