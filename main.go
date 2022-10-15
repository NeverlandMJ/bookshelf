package main

import (
	"github.com/NeverlandMJ/bookshelf/api"
	"github.com/NeverlandMJ/bookshelf/config"
	"github.com/NeverlandMJ/bookshelf/server"
	"github.com/NeverlandMJ/bookshelf/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	server, err := server.NewServer(cfg)
	if err != nil {
		panic(err)
	}

	defer server.DB.Close()

	service := service.NewService(server)
	if err != nil {
		panic(err)
	}

	r := api.NewRouter(service)

	r.Run(":" + cfg.Port)
}
