package models

import (
	"github.com/gitalek/gogal/hash"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	UserDB
}

func NewUserService(connStr string) (UserService, error) {
	ug, err := newUserGorm(connStr)
	if err != nil {
		return nil, err
	}
	hmac := hash.NewHMAC(hmacSecretKey)
	uv, err := newUserValidator(ug, hmac)
	if err != nil {
		return nil, err
	}
	return &userService{
		UserDB: uv,
	}, nil
}

func (us *userService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword(
		[]byte(foundUser.PasswordHash),
		[]byte(password+userPwPepper),
	)
	// Three use cases:
	switch err {
	case nil:
		return foundUser, nil
	case bcrypt.ErrMismatchedHashAndPassword:
		return nil, ErrPasswordIncorrect
	default:
		return nil, err
	}
}
