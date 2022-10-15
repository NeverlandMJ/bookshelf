package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/NeverlandMJ/bookshelf/pkg/customErr"
	"github.com/NeverlandMJ/bookshelf/pkg/entity"
	"github.com/NeverlandMJ/bookshelf/server"
)

// Service holds a type which implements all the methods of Repository.
// It acts as a middleman between api and database server
type Service struct {
	Repo server.Repository
}

// NewService creates a new Service
func NewService(repo server.Repository) *Service {
	return &Service{
		Repo: repo,
	}
}

func (s Service) SaveUser(ctx context.Context, u entity.UserSignUpRequest) (entity.ResponseUser, error) {
	user, err := s.Repo.SaveUser(ctx, u)
	if err != nil {
		return entity.ResponseUser{}, err
	}

	return entity.ConvertToResponseUser(user), nil
}

func (s Service) GetUser(ctx context.Context, key string) (entity.ResponseUser, error) {
	user, err := s.Repo.GetUser(ctx, key)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.ResponseUser{}, customErr.ErrNotFound
		}
		return entity.ResponseUser{}, err
	}

	return entity.ConvertToResponseUser(user), nil
}

func (s Service) SaveBook(ctx context.Context, book entity.Book) (entity.ResponseBook, error) {
	return s.Repo.SaveBook(ctx, book)
}

func (s Service) GetBook(ctx context.Context, isbn string) (entity.ResponseBook, error) {
	b, err := s.Repo.GetBook(ctx, isbn)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.ResponseBook{}, customErr.ErrNotFound
		}
		return entity.ResponseBook{}, err
	}

	return entity.ConvertToResponseBook(b), nil
}
func (s Service) GetAllBooks(ctx context.Context) ([]entity.ResponseBook, error) {
	books, err := s.Repo.GetAllBooks(ctx)
	if err != nil {
		return nil, err
	}

	respBooks := make([]entity.ResponseBook, 0, len(books))
	for _, book := range books {
		respBooks = append(respBooks, entity.ConvertToResponseBook(book))
	}

	return respBooks, nil
}

func (s Service) EditBook(ctx context.Context, status entity.EditBookReq, id int) (entity.ResponseBook, error) {
	book, err := s.Repo.EditBook(ctx, status, id)
	if err != nil {
		return entity.ResponseBook{}, err
	}

	return entity.ConvertToResponseBook(book), nil
}

func (s Service) DeleteBook(ctx context.Context, id int) error {
	return s.Repo.DeleteBook(ctx, id)
}
