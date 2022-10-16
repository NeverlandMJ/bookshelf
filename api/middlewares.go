package api

import (
	"context"
	"crypto/md5"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h Handler) Authentication(c *gin.Context) {
	key := c.GetHeader("Key")
	sign := c.GetHeader("Sign")
	if key == "" ||  sign== "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"isOk":    false,
			"message": "1user is unauthenticated",
		})
		c.Abort()
		return
	}

	got, err := h.srvc.GetUser(context.Background(), key)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"isOk":    false,
			"message": "1user is unauthenticated",
		})
		c.Abort()
		return
	}

	r := c.Request.URL
	url := r.String()

	body := []byte{}
	c.Request.Body.Read(body)

	secretByte := md5.Sum([]byte(c.Request.Method + url + string(body) + got.Secret))

	secret := fmt.Sprintf("%x", secretByte)

	if secret != sign {
		c.JSON(http.StatusUnauthorized, gin.H{
			"isOk":    false,
			"message": "1user is unauthenticated",
		})
		c.Abort()
		return
	}

	if secret != sign {
		c.JSON(http.StatusUnauthorized, gin.H{
			"isOk":    false,
			"message": "1user is unauthenticated",
		})
		c.Abort()
		return
	}

	c.Next()
}
