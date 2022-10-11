package api

import "github.com/NeverlandMJ/bookshelf/service"

type Handler struct {
	srvc *service.Service
}

func NewHandler(srv *service.Service) Handler {
	return Handler{
		srvc: srv,
	}
}