package api

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/NeverlandMJ/bookshelf/pkg/customErr"
	"github.com/gin-gonic/gin"
)

func (h Handler) Authentication(c *gin.Context) {
	key := c.GetHeader("Key")
	sign := c.GetHeader("Sign")
	if key == "" || sign == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"isOk":    false,
			"message": "user is unauthenticated: header is empty",
		})
		c.Abort()
		return
	}

	// check if the user's secret has been catched
	value, found := newCache.Get(key)
	userSecret, ok := value.(string)
	if !found || !ok {
		// if user's secret key has not been cached or due to crashes on server cach has been destroyed
		// we will get user key directly from database
		// this method is used to reduce number of database calls
		user, err := h.srvc.GetUser(context.Background(), key)
		if err != nil {
			if errors.Is(err, customErr.ErrNotFound) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"isOk":    false,
					"message": err.Error(),
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"isOk":    false,
				"message": err.Error(),
			})
			c.Abort()
			return
		}
		userSecret = user.Secret
	}
	scheme := "http://"
	if c.Request.TLS != nil {
		scheme = "https://"
	}
	url := scheme + c.Request.Host + c.Request.URL.Path

	jsonData, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Println(string(jsonData))

	secretByte := md5.Sum([]byte(c.Request.Method + url + string(jsonData) + userSecret))

	secret := fmt.Sprintf("%x", secretByte)

	if secret != sign {
		c.JSON(http.StatusUnauthorized, gin.H{
			"isOk":    false,
			"message": "incorrect sign",
		})
		c.Abort()
		return
	}

	c.Next()
}
