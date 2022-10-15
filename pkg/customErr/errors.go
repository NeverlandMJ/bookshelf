package customErr

import "errors"

var ErrInvalidInput = errors.New("invalid user input")
var ErrNotFound = errors.New("requested data is not found")
