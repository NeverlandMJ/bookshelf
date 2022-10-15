package server

import (
	"context"
	"github.com/NeverlandMJ/bookshelf/pkg/entity"

	"github.com/NeverlandMJ/bookshelf/config"
	"github.com/NeverlandMJ/bookshelf/database"
	"github.com/jmoiron/sqlx"
)

// Server holds databse
type Server struct {
	DB *sqlx.DB
}

// NewServer returns a new Server with working database attached to it.
// If an error occurs while connecting to database, it returns an error
func NewServer(cnfg config.Config) (*Server, error) {
	conn, err := database.Connect(cnfg)
	if err != nil {
		return nil, err
	}
	return &Server{
		DB: conn,
	}, nil
}

func (s Server) SaveUser(ctx context.Context, user entity.UserSignUpRequest) (entity.UserResponseFromDatabse, error) {
	query := `INSERT INTO users (name, key, secret) VALUES ($1, $2, $3) RETURNING id`

	var id int
	err := s.DB.QueryRowContext(ctx, query, user.Name, user.Key, user.Secret).Scan(&id)
	if err != nil {
		return entity.UserResponseFromDatabse{}, err
	}
	return entity.UserResponseFromDatabse{
		id,
		user.Name,
		user.Key,
		user.Secret,
	}, nil
}

func (s Server) GetUser(ctx context.Context, key string) (entity.UserResponseFromDatabse, error) {
	query := `SELECT id, name, key, secret FROM users WHERE key=$1`
	var user entity.UserResponseFromDatabse
	if err := s.DB.GetContext(ctx, &user, query, key); err != nil {
		return entity.UserResponseFromDatabse{}, err
	}

	return user, nil
}

func (s Server) SaveBook(ctx context.Context, book entity.Book) (entity.ResponseBook, error) {
	query := `INSERT INTO books (isbn, title, author, published, pages) 
	VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := s.DB.QueryRowContext(ctx, query,
		book.Isbn,
		book.Title,
		book.Author,
		book.Published,
		book.Pages,
	).Scan(&book.ID)
	if err != nil {
		return entity.ResponseBook{}, err
	}

	return entity.ResponseBook{
		book,
		0,
	}, nil
}

func (s Server) GetBook(ctx context.Context, isbn string) (entity.BookResponseFromDatabase, error) {
	query := `SELECT id, isbn, title, author, published, pages, status FROM books WHERE isbn=$1`
	var book entity.BookResponseFromDatabase
	if err := s.DB.GetContext(ctx, &book, query, isbn); err != nil {
		return entity.BookResponseFromDatabase{}, err
	}

	return book, nil
}

func (s Server) GetAllBooks(ctx context.Context) ([]entity.BookResponseFromDatabase, error) {
	query := `SELECT id, isbn, title, author, published, pages, status FROM books`
	books := make([]entity.BookResponseFromDatabase, 0)

	if err := s.DB.SelectContext(ctx, &books, query); err != nil {
		return nil, err
	}

	return books, nil
}

func (s Server) EditBook(ctx context.Context, status, id int) (entity.BookResponseFromDatabase, error) {
	query := `UPDATE books SET status=$1 WHERE id=$2 RETURNING 
		id, isbn, title, author, published, pages, status 
	`

	var book entity.BookResponseFromDatabase
	if err := s.DB.GetContext(ctx, &book, query, status, id); err != nil {
		return entity.BookResponseFromDatabase{}, err
	}

	return book, nil
}

func (s Server) DeleteBook(ctx context.Context, id int) error {
	query := `DELETE FROM books WHERE id=$1`

	_, err := s.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
