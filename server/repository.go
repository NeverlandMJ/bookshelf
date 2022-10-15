package server

import (
	"context"
	"github.com/NeverlandMJ/bookshelf/pkg/entity"
)

type Repository interface {
	SaveUser(ctx context.Context, user entity.UserSignUpRequest) (entity.UserResponseFromDatabse, error)
	GetUser(ctx context.Context, key string) (entity.UserResponseFromDatabse, error)
	SaveBook(ctx context.Context, book entity.Book) (entity.ResponseBook, error)
	GetBook(ctx context.Context, isbn string) (entity.BookResponseFromDatabase, error)
	GetAllBooks(ctx context.Context) ([]entity.BookResponseFromDatabase, error)
	EditBook(ctx context.Context, status, id int) (entity.BookResponseFromDatabase, error)
	DeleteBook(ctx context.Context, id int) error
}
