package models

import (
	"github.com/gitalek/gogal/hash"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	UserDB
	pepper string
}

func NewUserService(db *gorm.DB, pepper, hmacKey string) UserService {
	ug := &userGorm{db: db}
	hmac := hash.NewHMAC(hmacKey)
	uv, err := newUserValidator(ug, hmac, pepper)
	if err != nil {
		return nil
	}
	return &userService{
		UserDB: uv,
		pepper: pepper,
	}
}

func (us *userService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword(
		[]byte(foundUser.PasswordHash),
		[]byte(password+us.pepper),
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
