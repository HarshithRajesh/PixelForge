package processor

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Health Check APIs are working",
	})
}

func Profile(c *gin.Context) {
	email := c.MustGet("userID").(string)
	c.JSON(http.StatusOK, gin.H{"message": "Protected route check", "email": email})
}
