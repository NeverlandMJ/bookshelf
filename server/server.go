package server

import (
	"github.com/NeverlandMJ/bookshelf/config"
	"github.com/NeverlandMJ/bookshelf/database"
	"github.com/jmoiron/sqlx"
)

// Server holds databse
type Server struct {
	DB *sqlx.DB
}

// NewServer returns a new Server with working database attached to it.
// If an error occuras while connecting to database, it returns an error
func NewServer(cnfg config.Config) (*Server, error) {
	conn, err := database.Connect(cnfg)
	if err != nil {
		return nil, err
	}
	return &Server{
		DB: conn,
	}, nil
}


