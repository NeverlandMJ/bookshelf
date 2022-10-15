package server

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/NeverlandMJ/bookshelf/config"
	"github.com/NeverlandMJ/bookshelf/pkg/customErr"
	"github.com/NeverlandMJ/bookshelf/pkg/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testUser1 = entity.UserSignUpRequest{
	Name:   "Jackson",
	Key:    "MyKey",
	Secret: "MySecret",
}

var testUser2 = entity.UserSignUpRequest{
	Name:   "Amanda",
	Key:    "MyKey",
	Secret: "secret",
}

var testBook1 = entity.Book{
	Isbn:      "9781118464465",
	Title:     "Raspberry Pi User Guide",
	Author:    "Eben Upton",
	Published: 2012,
	Pages:     221,
}

var testBook2 = entity.Book{
	Isbn:      "9874563217896",
	Title:     "Paper Towns",
	Author:    "JHon Green",
	Published: 2012,
	Pages:     221,
}

func TestServer_SaveUser(t *testing.T) {
	s := newServer(t)
	t.Cleanup(cleanUpFn(s))

	t.Run("success", func(t *testing.T) {
		t.Cleanup(cleanUpFn(s))

		want := entity.UserResponseFromDatabse{
			ID:     1,
			Name:   "Jackson",
			Key:    "MyKey",
			Secret: "MySecret",
		}

		got, err := s.SaveUser(context.Background(), testUser1)
		require.NoError(t, err)
		assert.EqualValues(t, want, got)

	})

	t.Run("duplicate user key", func(t *testing.T) {
		t.Cleanup(cleanUpFn(s))

		_, err := s.SaveUser(context.Background(), testUser1)
		require.NoError(t, err)

		got, err := s.SaveUser(context.Background(), testUser2)
		require.ErrorIs(t, err, customErr.ErrSaveUserError)
		assert.Empty(t, got)
	})
}

func TestServer_GetUser(t *testing.T) {
	s := newServer(t)
	t.Cleanup(cleanUpFn(s))

	t.Run("success", func(t *testing.T) {
		t.Cleanup(cleanUpFn(s))

		want := entity.UserResponseFromDatabse{
			ID:     1,
			Name:   "Jackson",
			Key:    "MyKey",
			Secret: "MySecret",
		}

		_, err := s.SaveUser(context.Background(), testUser1)
		require.NoError(t, err)

		got, err := s.GetUser(context.Background(), "MyKey")
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("user doesn't exist", func(t *testing.T) {
		t.Cleanup(cleanUpFn(s))

		got, err := s.GetUser(context.Background(), "MyKey")
		require.ErrorIs(t, err, sql.ErrNoRows)
		assert.Empty(t, got)
	})
}

func TestServer_SaveBooks(t *testing.T) {
	s := newServer(t)
	t.Cleanup(cleanUpFn(s))

	t.Run("sucess", func(t *testing.T) {
		t.Cleanup(cleanUpFn(s))

		want := entity.ResponseBook{
			Book: entity.Book{
				ID:        1,
				Isbn:      "9781118464465",
				Title:     "Raspberry Pi User Guide",
				Author:    "Eben Upton",
				Published: 2012,
				Pages:     221,
			},
			Status: 0,
		}

		got, err := s.SaveBook(context.Background(), testBook1)
		require.NoError(t, err)
		assert.EqualValues(t, want, got)
	})
}

func TestServer_GetBook(t *testing.T) {
	s := newServer(t)
	t.Cleanup(cleanUpFn(s))

	t.Run("success", func(t *testing.T) {
		t.Cleanup(cleanUpFn(s))

		want := entity.BookResponseFromDatabase{
			ID:        1,
			Isbn:      "9781118464465",
			Title:     "Raspberry Pi User Guide",
			Author:    "Eben Upton",
			Published: 2012,
			Pages:     221,
			Status: 0,
		}

		_, err := s.SaveBook(context.Background(), testBook1)
		require.NoError(t, err)

		got, err := s.GetBook(context.Background(), want.Isbn)
		require.NoError(t, err)
		assert.EqualValues(t, want, got)

	})

	t.Run("book doesn't exist", func(t *testing.T) {
		t.Cleanup(cleanUpFn(s))

		got, err := s.GetBook(context.Background(), "9781118464465")
		require.ErrorIs(t, err, sql.ErrNoRows)
		assert.Empty(t, got)

	})
}

func TestServer_GetAllBooks(t *testing.T) {
	s := newServer(t)
	t.Cleanup(cleanUpFn(s))

	t.Run("success", func(t *testing.T) {
		t.Cleanup(cleanUpFn(s))

		_, err := s.SaveBook(context.Background(), testBook1)
		require.NoError(t, err)

		_, err = s.SaveBook(context.Background(), testBook2)
		require.NoError(t, err)

		b1, err := s.GetBook(context.Background(), testBook1.Isbn)
		require.NoError(t, err)

		b2, err := s.GetBook(context.Background(), testBook2.Isbn)
		require.NoError(t, err)

		want := []entity.BookResponseFromDatabase{
			b1, b2,
		}

		got, err := s.GetAllBooks(context.Background())
		require.NoError(t, err)
		assert.EqualValues(t, want, got)

	})
}


func TestServer_EditBook(t *testing.T) {
	s := newServer(t)
	t.Cleanup(cleanUpFn(s))

	t.Run("sucess", func(t *testing.T) {
		t.Cleanup(cleanUpFn(s))

		b, err := s.SaveBook(context.Background(), testBook1)
		require.NoError(t, err)

		got, err := s.EditBook(context.Background(), entity.EditBookReq{
			Book: b.Book,
			Status: 1,
		}, b.Book.ID)

		require.NoError(t, err)

		want, err := s.GetBook(context.Background(), got.Isbn)
		require.NoError(t, err)

		assert.Equal(t, want, got)
	})


}

func TestServer_DeleteBook(t *testing.T) {
	s := newServer(t)
	t.Cleanup(cleanUpFn(s))

	t.Run("success", func(t *testing.T) {
		t.Cleanup(cleanUpFn(s))

		b, err := s.SaveBook(context.Background(), testBook1)
		require.NoError(t, err)

		err = s.DeleteBook(context.Background(), b.Book.ID)
		require.NoError(t, err)

		want, err := s.GetBook(context.Background(), b.Book.Isbn)
		require.ErrorIs(t, err, sql.ErrNoRows)
		assert.Empty(t, want)
	})
}

func newServer(t *testing.T) *Server {
	t.Helper()
	serv, err := NewServer(
		config.Config{
			Host:                   "localhost",
			Port:                   "8080",
			PostgresHost:           "localhost",
			PostgresPort:           "5432",
			PostgresUser:           "postgres",
			PostgresPassword:       "1234",
			PostgresDB:             "postgres",
			PostgresMigrationsPath: "file://./../../migrations",
		},
	)

	require.NoError(t, err)

	return serv
}

func cleanUpFn(s *Server) func() {
	return func() {
		if err := s.delete(); err != nil {
			log.Println("CLEANUP OF DB FAILED!", err.Error())
		}
	}
}
