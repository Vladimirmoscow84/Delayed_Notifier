package handlers

import (
	"fmt"
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

	id, err := r.service.SaveData(ctx, notice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sd := notice.SendDate.String()
	fmt.Printf("sd: %s\n", sd)
	t, err := time.Parse("2006-01-02 15:04:05 -0700 UTC", sd)
	if err != nil {
		fmt.Printf("error parsing time: %v\n", err)
	} else {
		fmt.Printf("t: %v\n", t)
	}

	notice.Id = id

	fmt.Printf("notice: %v\n", notice)

	err = r.store.Cache.Set(ctx, sd, &notice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newNotice, err := r.store.Cache.Get(ctx, sd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("newNotice: %v\n", *newNotice)

	c.JSON(http.StatusOK, gin.H{"id": id})
}
