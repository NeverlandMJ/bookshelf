package api

import (
	"fmt"

	"github.com/NeverlandMJ/bookshelf/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(serv *service.Service) *gin.Engine {
	router := gin.Default()
	h := NewHandler(serv)

	// solving cors error
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	router.Use(cors.New(config))

	router.POST("/signup", h.SignUp)
	router.GET("/hi", func(ctx *gin.Context) {
		fmt.Fprintln(ctx.Writer, "App is running")
	})

	authorized := router.Group("/")
	authorized.Use(Authentication)
	authorized.GET("/myself", h.GetUser)
	authorized.POST("/books", h.SaveBook)
	authorized.GET("/books", h.GetAllBooks)
	authorized.PATCH("/books/:id", h.EditBook)
	authorized.DELETE("/books/:id", h.DeleteBook)
	return router
}
