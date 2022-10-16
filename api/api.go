package api

import (
	"fmt"

	"github.com/NeverlandMJ/bookshelf/service"
	"github.com/gin-gonic/gin"
)

func NewRouter(serv *service.Service) *gin.Engine {
	router := gin.Default()
	h := NewHandler(serv)

	router.POST("/signup", h.SignUp)
	router.GET("/cleanup", h.Delete)

	router.GET("/", func(ctx *gin.Context) {
		fmt.Fprintln(ctx.Writer, "App is running")
	})

	auth := router.Group("/")
	auth.Use(h.Authentication)
	auth.GET("/myself", h.GetUser)
	auth.POST("/books", h.SaveBook)
	auth.GET("/books", h.GetAllBooks)
	auth.PATCH("/books/:id", h.EditBook)
	auth.DELETE("/books/:id", h.DeleteBook)

	return router
}
