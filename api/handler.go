package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/NeverlandMJ/bookshelf/pkg/customErr"
	"github.com/NeverlandMJ/bookshelf/pkg/entity"
	"github.com/NeverlandMJ/bookshelf/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	srvc *service.Service
}

func NewHandler(srv *service.Service) Handler {
	return Handler{
		srvc: srv,
	}
}

func (h Handler) SignUp(c *gin.Context) {
	var req entity.UserSignUpRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"isOk":    false,
			"message": err.Error(),
		})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"isOk":    false,
			"message": err.Error(),
		})
		return
	}

	user, err := h.srvc.SaveUser(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"isOk":    false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, entity.Response{
		Data:    user,
		IsOk:    true,
		Message: "ok",
	})
}

func (h Handler) GetUser(c *gin.Context) {
	key := c.GetHeader("Key")

	user, err := h.srvc.GetUser(context.Background(), key)
	if err != nil {
		// if errors.Is(err, customErr.ErrNotFound) {
		// 	c.JSON(http.StatusNotFound, gin.H{
		// 		"isOk":    false,
		// 		"message": err.Error(),
		// 	})
		// 	return
		// }
		c.JSON(http.StatusInternalServerError, gin.H{
			"isOk":    false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, entity.Response{
		Data:    user,
		IsOk:    true,
		Message: "ok",
	})
}

func (h Handler) SaveBook(c *gin.Context) {
	// reading isbn from body
	var isbn entity.CreatBookRequest
	if err := c.BindJSON(&isbn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"isOk":    false,
			"message": err.Error(),
		})
		return
	}
	// checking if isbn valid
	if err := isbn.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"isOk":    false,
			"message": err.Error(),
		})
		return
	}

	// checking if the book has already been saved
	book, err := h.srvc.GetBook(context.Background(), isbn.Isbn)
	if err != nil {
		// if it was not saved before then we will send request to Open API to get book info
		if errors.Is(err, customErr.ErrNotFound) {
			bookInfo, err := getBook(isbn.Isbn)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"isOk":    false,
					"message": err.Error(),
				})
				return
			}
			// after fetching book info from Open API we will save it to the database
			book, err := h.srvc.SaveBook(context.Background(), bookInfo)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"isOk":    false,
					"message": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, entity.Response{
				Data:    book,
				IsOk:    true,
				Message: "ok",
			})
			return
		}

		// if there was any internal error while getting book from database return error
		c.JSON(http.StatusInternalServerError, gin.H{
			"isOk":    false,
			"message": err.Error(),
		})
		return
	}

	// if the book has already been fetched and saved to database
	// just get that book from db and return to the user
	c.JSON(http.StatusOK, entity.Response{
		Data:    book,
		IsOk:    true,
		Message: "ok",
	})

}

func (h Handler) GetAllBooks(c *gin.Context) {
	books, err := h.srvc.GetAllBooks(context.Background())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"isOk":    false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, entity.Response{
		Data:    books,
		IsOk:    true,
		Message: "ok",
	})
}

func (h Handler) EditBook(c *gin.Context) {
	pathID, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusBadRequest, entity.Response{
			IsOk:    false,
			Message: "ID is not provided",
		})
		return
	}

	id, err := strconv.Atoi(pathID)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Response{
			IsOk:    false,
			Message: "invalid ID",
		})
	}

	var status entity.EditBookReq
	if err := c.BindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, entity.Response{
			IsOk:    false,
			Message: err.Error(),
		})
		return
	}

	book, err := h.srvc.EditBook(context.Background(), status, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Response{
			IsOk:    false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, entity.Response{
		Data:    book,
		IsOk:    true,
		Message: "ok",
	})
}

func (h Handler) DeleteBook(c *gin.Context) {
	pathID, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusBadRequest, entity.Response{
			IsOk:    false,
			Message: "ID is not provided",
		})
		return
	}

	id, err := strconv.Atoi(pathID)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Response{
			IsOk:    false,
			Message: "invalid ID",
		})
	}

	err = h.srvc.DeleteBook(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Response{
			IsOk:    false,
			Message: err.Error(),
		})
	}

	c.JSON(http.StatusOK, entity.Response{
		Data:    "Successfully deleted",
		IsOk:    true,
		Message: "ok",
	})
}

func getBook(isbn string) (entity.Book, error) {
	resp, err := http.Get("https://openlibrary.org/isbn/" + url.QueryEscape(isbn) + ".json")
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var info entity.Info
	if err := json.Unmarshal(body, &info); err != nil {
		fmt.Println(err)
	}
	pd := strings.Fields(info.PublishDate)
	date, err := strconv.Atoi(pd[len(pd)-1])
	if err != nil {
		return entity.Book{}, err
	}
	if len(info.Authors) <= 0 {
		return entity.Book{
			Isbn:      isbn,
			Title:     info.Title,
			Author:    "",
			Published: date,
			Pages:     info.NumberOfPages,
		}, nil
	}
	resp, err = http.Get("https://openlibrary.org/" + url.QueryEscape(info.Authors[0].Key) + ".json")
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var author entity.Author
	if err := json.Unmarshal(body, &author); err != nil {
		fmt.Println(err)
	}

	return entity.Book{
		Isbn:      isbn,
		Title:     info.Title,
		Author:    author.PersonalName,
		Published: date,
		Pages:     info.NumberOfPages,
	}, nil
}
