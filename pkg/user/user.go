package user

import (
	"errors"
)

type User struct {
	UserName string
	Password string
}

// ValidateUser - given a user validate its field values.
func (u User) ValidateUser() error {
	if u.UserName == "" || u.Password == "" {
		return errors.New("invalid user data, empty fields")
	}
	return nil
}
