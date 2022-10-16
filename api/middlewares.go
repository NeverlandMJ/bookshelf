package api

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authentication(c *gin.Context) {
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

	userSecret, found := newCache.Get(key)
	if !found {
		c.JSON(http.StatusUnauthorized, gin.H{
			"isOk":    false,
			"message": "user is unauthenticated: session is empty",
		})
		c.Abort()
		return
	} else {
		scheme := "https://"
		// if c.Request.TLS == nil {
		// 	scheme = "http://"
		// }
		url := scheme + c.Request.Host + c.Request.URL.Path

		jsonData, _ := ioutil.ReadAll(c.Request.Body)
		fmt.Println(string(jsonData))

		secretByte := md5.Sum([]byte(c.Request.Method + url + string(jsonData) + userSecret.(string)))

		secret := fmt.Sprintf("%x", secretByte)

		fmt.Println("method: ", c.Request.Method)
		fmt.Println("url: ", url)
		fmt.Println("body: ", string(jsonData))
		fmt.Println("secret: ", userSecret)

		if secret != sign {
			c.JSON(http.StatusUnauthorized, gin.H{
				"isOk":    false,
				"message": "incorrect sign",
			})
			c.Abort()
			return
		}
	}

	c.Next()
}
