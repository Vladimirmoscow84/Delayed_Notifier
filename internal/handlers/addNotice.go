package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/model"
	"github.com/gin-gonic/gin"
)

// Структура уведомления, отправляемая в RabbitMQ.
type NotifyMessage struct {
	Subject string    `json:"subject"`
	Body    string    `json:"body"`
	SendAt  time.Time `json:"send_at"`
}

func (r *Router) addNotice(c *gin.Context) {
	ctx := c.Request.Context()

	var notice model.Notice

	if err := c.ShouldBindJSON(&notice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input data."})
		return
	}

	now := time.Now()
	notice.DateCreated = now
	notice.SendAttempts = 5
	notice.SendStatus = "sheduled"
	fmt.Println("[REDIS]Отправка в базу")
	err := r.dataSaver.SaveData(ctx, notice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("[REDIS]Отправка в базу завершена")

	if r.rabbit != nil {
		delay := time.Until(notice.SendDate)
		if delay < 0 {
			delay = 0
		}
		fmt.Println("[RABBITMQ]Отправка на рассылку")

		if r.rabbit != nil {
			fmt.Printf("[RABBITMQ] Публикация уведомления с задержкой %v...\n", delay)

			subject := fmt.Sprintf("Уведомление №%d", notice.Id)
			msg := NotifyMessage{
				Subject: subject,
				Body:    notice.Body,
				SendAt:  notice.SendDate,
			}

			err := r.rabbit.PublishStructWithTTL(msg, delay)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "failed to publish message to RabbitMQ: " + err.Error(),
				})
				return
			}

			fmt.Println("[RABBITMQ] Публикация успешна.")
		}

		c.JSON(http.StatusOK, gin.H{
			"id":      strconv.Itoa(notice.Id),
			"status":  "scheduled",
			"send_at": notice.SendDate.Format(time.RFC3339),
		})
	}
}
