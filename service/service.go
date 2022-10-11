package service

import (
	"github.com/NeverlandMJ/bookshelf/server"
)

// Service holds a type which implements all the methods of Repository.
// It acts as a middleman between api and database server
type Service struct {
	Repo server.Repository
}

// NewService creates a new Service
func NewService(repo server.Repository, redisAddr string) *Service {
	return &Service{
		Repo: repo,
	}
}


