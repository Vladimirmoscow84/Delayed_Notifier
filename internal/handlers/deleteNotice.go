package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (r *Router) deleteNotice(c *gin.Context) {
	ctx := c.Request.Context()
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing ID parameter of notice"})
		return
	}

	err = r.dataDeleter.DeleteData(ctx, idStr)
	if err != nil {
		log.Printf("error deleting notice by ID=%d: %v\n", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "an internal error while deleting the notice by ID"})
		return
	}
	c.JSON(http.StatusOK, gin.H{idStr: "notice has successfully deleted by ID from redis"})
}
