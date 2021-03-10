package models

import "github.com/jinzhu/gorm"

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
	// Close is used to close a DB connection.
	Close() error
	// Migration helpers
	AutoMigrate() error
	DestructiveReset() error
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

