package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	// ErrNotFound is returned when a resource cannot be found in the database
	ErrNotFound = errors.New("models: resource not found")
)

type User struct {
	gorm.Model
	Name string
	Email string `gorm:"not null;unique_index"`
}

type UserService struct {
	db *gorm.DB
}

func NewUserService(connStr string) (*UserService, error) {
	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &UserService{db: db}, nil
}

// Close method closes the UserService database connection.
func (us *UserService) Close() error {
	return us.db.Close()
}

// DestructiveReset method drops the user table and rebuilds it.
// Used in development env.
func (us *UserService) DestructiveReset() {
	us.db.DropTableIfExists(&User{})
	us.db.AutoMigrate(&User{})
}

func (us *UserService) Create(user *User) error {
	return us.db.Create(user).Error
}

// As a general rule, any error but ErrNotFound should
// probably result in a 500 error.
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	err := us.db.Where("id = ?", id).First(&user).Error
	switch err {
	case nil:
		return &user, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
