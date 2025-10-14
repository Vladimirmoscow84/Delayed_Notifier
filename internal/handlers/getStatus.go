package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (r *Router) getStatus(c *gin.Context) {
	ctx := c.Request.Context()
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil && id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing ID parametr of notice"})
		return
	}
	statusNotice, err := r.statusGetter.GetStatusNotice(ctx, idStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"sendStatus": statusNotice})
}
