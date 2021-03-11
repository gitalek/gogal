package models

import (
	"github.com/jinzhu/gorm"
	"strings"
)

var (
	// ErrNotFound is returned when a resource cannot be found in the database.
	ErrNotFound modelError = "models: resource not found"
	// ErrIdInvalid is returned when an invalid ID is provided to a method like Delete.
	ErrIdInvalid modelError = "models: ID provided was invalid"
	// ErrPasswordIncorrect is returned when an invalid password is used when attempting
	// to authenticate a user.
	ErrPasswordIncorrect modelError = "models: incorrect password provided"
	// ErrEmailRequired is returned when an email address is not provided when creating a user.
	ErrEmailRequired modelError = "models: email address is required"
	// ErrEmailInvalid is returned when an email address doesn't match regexp.
	ErrEmailInvalid modelError = "models: email address is not valid"
	// ErrEmailTaken is returned when an update or create is attempted with an email address
	// that is already in use.
	ErrEmailTaken modelError = "models: email address is already taken"
	// ErrPasswordTooShort is returned when a user tries to set a password
	// that is less than 8 characters long.
	ErrPasswordTooShort modelError = "models: password must be at least 8 characters long"
	// ErrPasswordRequired is returned when a create is attempted without a user password provided.
	ErrPasswordRequired modelError = "models: password is required"
	// ErrRememberRequired is returned when a create or update
	// is attempted without a user remember token hash.
	ErrRememberRequired modelError = "models: remember token is required"
	// ErrRememberTooShort is returned when a remember token is not at least 32 bytes.
	ErrRememberTooShort modelError = "models: remember token must be at least 32 bytes"
	// userPwPepper is used for peppering passwords.
	userPwPepper = "secret-random-string"
	// hmacSecretKey is used for hashing remember tokens.
	hmacSecretKey = "secret-hmac-key"
)

// UserDB is used to interact with the users database.
type UserDB interface {
	// Query for single users.
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)
	// Alter users.
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error
}

// UserService is a set of methods used to manipulate and work with the user model
type UserService interface {
	UserDB
	// Authenticate can be used to authenticate a user with the
	// provided email address and password.
	// If the email address provided is invalid, this will return
	//   nil, ErrNotFound
	// If the password provided is invalid, this will return
	//   nil, ErrPasswordIncorrect
	// If the email and password are both valid, this will return
	//   user, nil
	// Otherwise if another error is encountered this will return
	//   nil, error
	Authenticate(email, password string) (*User, error)
}

type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}
