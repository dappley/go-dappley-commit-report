package helper

import (
	"errors"
)

//Return an error message when the input flag argument is a default value.
func CheckFlags(email string, password string) (err error) {
	switch {
	case email == "default_email":
		err = errors.New("Error: Email is missing!")
	case password == "default_password":
		err = errors.New("Error: Password is missing!")
	default:
		err = nil
	}
	return
}