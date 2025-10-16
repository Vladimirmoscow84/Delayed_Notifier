package telegram

import (
	"errors"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Client struct {
	chatID int64
	bot    *tgbotapi.BotAPI
}

func New(token string, chatID int64) (*Client, error) {
	if token == "" {
		log.Println("[telega]TELEGRAM_TOKEN is not specified, no messages will be sent in telegram")
		return nil, errors.New("[telega]TELEGRAM_TOKEN is not specified")
	}
	if chatID == 0 {
		log.Println("[telega]TELEGRAM_CHAT_ID is not specified, no messages will be sent in telegram")
		return nil, errors.New("[telega]TELEGRAM_CHAT_ID is not specified")
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	log.Printf("[telega] Client successfully created. Authorized on account %s", bot.Self.UserName)
	return &Client{
		chatID: chatID,
		bot:    bot,
	}, nil
}

// Send - отправляет уведомление в телеграм
func (c *Client) Send(message string) error {
	if c.bot == nil {
		return errors.New("[telega] bot client is nil")
	}
	msg := tgbotapi.NewMessage(c.chatID, message)
	msg.DisableWebPagePreview = true
	_, err := c.bot.Send(msg)
	if err != nil {
		return err
	}
	log.Printf("[telega] message successfully sent to chatID %d", c.chatID)
	return nil
}
