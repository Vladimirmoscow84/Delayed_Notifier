package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (r *Router) getStatus(c *gin.Context) {
	ctx := c.Request.Context()
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing ID parametr of notice"})
		return
	}
	notice, err := r.store.GetNotice(ctx, id)
	if err != nil {
		log.Printf("error fetching notice by ID = %d: %v\n", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed getting status notice by ID"})
		return
	}
	c.JSON(http.StatusOK, notice.SendStatus)
}
