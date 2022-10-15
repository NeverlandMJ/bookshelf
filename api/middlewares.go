package api

import (
	"crypto/md5"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authentication(c *gin.Context) {
	key := c.GetHeader("Key")
	secret := c.GetHeader("Sign")
	if key == "" || secret == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"isOk":    false,
			"message": "user is unauthenticated",
		})
		c.Abort()
		return
	}

	session, _ := Store.Get(c.Request, key)
	var authenticated interface{} = session.Values[key]
	if authenticated != nil {
		authSecreteKey := session.Values[key].(string)
		if authSecreteKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"isOk":    false,
				"message": "user is unauthenticated",
			})
			c.Abort()
			return
		} else {
			body := []byte{}
			c.Request.Body.Read(body)
			sign := md5.Sum([]byte(c.Request.Method + c.Request.URL.String() + string(body) + authSecreteKey))
			signStr := fmt.Sprintf("%x", sign)
			if signStr != secret {
				c.JSON(http.StatusUnauthorized, gin.H{
					"isOk":    false,
					"message": "user is unauthenticated",
				})
				c.Abort()
				return
			}
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"isOk":    false,
			"message": "user is unauthenticated",
		})
		c.Abort()
		return
	}
	c.Next()
}
