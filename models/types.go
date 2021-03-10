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

type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}

