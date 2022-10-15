package entity

import (
	"fmt"
	"github.com/NeverlandMJ/bookshelf/pkg/customErr"
)

func (u UserSignUpRequest) Validate() error {
	if u.Name == "" {
		return fmt.Errorf("%w: empty user name", customErr.ErrInvalidInput)
	}
	if u.Key == "" {
		return fmt.Errorf("%w: empty key", customErr.ErrInvalidInput)
	}
	if u.Secret == "" {
		return fmt.Errorf("%w: empty secret", customErr.ErrInvalidInput)
	}
	return nil
}

func (c CreatBookRequest) Validate() error {
	if c.Isbn == "" {
		return fmt.Errorf("%w: empty isbn", customErr.ErrInvalidInput)
	}
	return nil
}
