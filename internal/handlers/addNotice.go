package handlers

import (
	"net/http"
	"time"

	"github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/model"
	"github.com/gin-gonic/gin"
)

func (r *Router) addNotice(c *gin.Context) {
	ctx := c.Request.Context()

	var notice model.Notice

	if err := c.ShouldBindJSON(&notice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invlid input data."})
		return
	}

	now := time.Now()
	notice.DateCreated = now
	notice.SendAttempts = 5
	notice.SendStatus = "sheduled"

	id, err := r.dataSaver.SaveData(ctx, notice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if r.rabbit != nil {
		delay := time.Until(notice.SendDate)
		if delay < 0 {
			delay = 0
		}
		notice.Id = id
		if err := r.rabbit.PublishStruct(notice, delay); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to publish message to RabbitMQ: " + err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}
