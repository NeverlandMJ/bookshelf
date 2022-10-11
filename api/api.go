package api

import (
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

	

	return router
}
