package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authentication(c *gin.Context) {
	key := c.GetHeader("Key")
	sign := c.GetHeader("Sign")
	if key == "" || sign == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"isOk":    false,
			"message": "user is unauthenticated",
		})
		c.Abort()
		return
	}

	c.Next()
}
