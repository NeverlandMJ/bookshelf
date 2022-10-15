package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authentication(c *gin.Context) {
	key := c.GetHeader("Key")
	secret := c.GetHeader("Sign")
	if key == "" || secret == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"isOk":    false,
			"message": "1user is unauthenticated",
		})
		c.Abort()
		return
	}

	
	c.Next()
}
