package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Client struct {
	chatID int64
	bot    *tgbotapi.BotAPI
}

func New(token string, chatID int64) (*Client, error) {
	if token == "" || chatID == 0 {
		log.Println("[telega]TELEGRAM_TIKEN or TELEGRAM_CHAT_ID is not specified, no messages will be sent in telegram")
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	log.Println("[telega] Client successfully created")
	return &Client{
		chatID: chatID,
		bot:    bot,
	}, nil
}

// Send - отправляет уведомление в телеграм
func (c *Client) Send(message string) error {
	msg := tgbotapi.NewMessage(c.chatID, message)
	_, err := c.bot.Send(msg)
	if err != nil {
		return err
	}
	log.Println("[telega] message successfully sended")
	return nil
}
