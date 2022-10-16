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

// SaveUser inserts a new user into database. User key must be unique if it was not unique
// function returns customErr.ErrSaveUserError
func (s Server) SaveUser(ctx context.Context, user entity.UserSignUpRequest) (entity.UserResponseFromDatabse, error) {
	query := `INSERT INTO users (name, key, secret) VALUES ($1, $2, $3) RETURNING id`

	var id int
	err := s.DB.QueryRowContext(ctx, query, user.Name, user.Key, user.Secret).Scan(&id)
	if err != nil {
		return entity.UserResponseFromDatabse{}, err
	}
	return entity.UserResponseFromDatabse{
		ID:     id,
		Name:   user.Name,
		Key:    user.Key,
		Secret: user.Secret,
	}, nil
}

// GetUser fetches user data from database if exists.
// If the given user doesn't exist it retursn sql.ErrNoRows
func (s Server) GetUser(ctx context.Context, key string) (entity.UserResponseFromDatabse, error) {
	query := `SELECT id, name, key, secret FROM users WHERE key=$1`
	var user entity.UserResponseFromDatabse
	if err := s.DB.GetContext(ctx, &user, query, key); err != nil {
		return entity.UserResponseFromDatabse{}, err
	}

	return user, nil
}

// SaveBook inserts book info into databse with default status 0-new
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
		Book:   book,
		Status: 0,
	}, nil
}

// GetBook fetches book info from database by book's isbn.
// If the asked b[ks doesn't exist it returns sql.ErrNoRows error
func (s Server) GetBook(ctx context.Context, isbn string) (entity.BookResponseFromDatabase, error) {
	query := `SELECT id, isbn, title, author, published, pages, status FROM books WHERE isbn=$1`
	var book entity.BookResponseFromDatabase
	if err := s.DB.GetContext(ctx, &book, query, isbn); err != nil {
		return entity.BookResponseFromDatabase{}, err
	}

	return book, nil
}

// GetAllBooks fetches all books from database
func (s Server) GetAllBooks(ctx context.Context) ([]entity.BookResponseFromDatabase, error) {
	query := `SELECT id, isbn, title, author, published, pages, status FROM books`
	books := make([]entity.BookResponseFromDatabase, 0)

	if err := s.DB.SelectContext(ctx, &books, query); err != nil {
		return nil, err
	}

	return books, nil
}

// EditBook edits book's data
func (s Server) EditBook(ctx context.Context, status entity.EditBookReq, id int) (entity.BookResponseFromDatabase, error) {
	query := `UPDATE books SET 
		isbn=$1, title=$2, author=$3, published=$4, pages=$5, status=$6  
		WHERE id=$7 RETURNING 
		id, isbn, title, author, published, pages, status 
	`

	var book entity.BookResponseFromDatabase
	if err := s.DB.GetContext(ctx, &book, query,
		status.Book.Isbn,
		status.Book.Title,
		status.Book.Author,
		status.Book.Published,
		status.Book.Pages,
		status.Status,
		id); err != nil {
		return entity.BookResponseFromDatabase{}, err
	}

	return book, nil
}

// DeleteBook deletes book from database by id
func (s Server) DeleteBook(ctx context.Context, id int) error {
	query := `DELETE FROM books WHERE id=$1`

	_, err := s.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

// delete this functions was written to be used as a cleanUP function inside intigrations tests
func (s Server) delete() error {
	_, err := s.DB.Exec(`
		TRUNCATE users RESTART IDENTITY CASCADE;
	`)
	if err != nil {
		return err
	}

	_, err = s.DB.Exec(`
		TRUNCATE books RESTART IDENTITY CASCADE;
	`)
	if err != nil {
		return err
	}

	return nil
}
