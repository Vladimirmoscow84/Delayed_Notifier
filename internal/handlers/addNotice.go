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

	notice.DateCreated = time.Now()
	notice.SendAttempts = 5
	notice.SendStatus = "sheduled"

	id, err := r.dataSaver.SaveData(ctx, notice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}
