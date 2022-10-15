package entity

func ConvertToResponseBook(book BookResponseFromDatabase) ResponseBook {
	return ResponseBook{
		Book: Book{
			ID:        book.ID,
			Isbn:      book.Isbn,
			Title:     book.Title,
			Author:    book.Author,
			Published: book.Published,
			Pages:     book.Pages,
		},
		Status: book.Status,
	}
}

func ConvertToResponseUser(user UserResponseFromDatabse) ResponseUser {
	return ResponseUser{
		ID:     user.ID,
		Name:   user.Name,
		Key:    user.Key,
		Secret: user.Secret,
	}
}
