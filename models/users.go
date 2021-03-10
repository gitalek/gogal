package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	// ErrNotFound is returned when a resource cannot be found in the database.
	ErrNotFound = errors.New("models: resource not found")
	// ErrInvalidID is returned when an invalid ID is provided to a method like Delete.
	ErrInvalidID = errors.New("models: ID provided was invalid")
	// ErrInvalidPassword is reterned when an invalid password is used when attempting
	// to authenticate a user.
	ErrInvalidPassword = errors.New("models: incorrect password provided")
	// ErrEmailRequired is returned when an email address is not provided when creating a user.
	ErrEmailRequired = errors.New("models: email address is required")
	// ErrEmailInvalid is returned when an email address doesn't match regexp.
	ErrEmailInvalid = errors.New("models: email address is not valid")
	// ErrEmailTaken is returned when an update or create is attempted with an email address
	// that is already in use.
	ErrEmailTaken = errors.New("models: email address is already taken")
	// userPwPepper is used for peppering passwords.
	userPwPepper = "secret-random-string"
	// hmacSecretKey is used for hashing remember tokens.
	hmacSecretKey = "secret-hmac-key"
)

// userGorm represents database interaction layer and implements the UserDB interface fully.
type userGorm struct {
	db   *gorm.DB
}

// Check if userGorm type implements UserDB interface.
var _ UserDB = &userGorm{}

func newUserGorm(connStr string) (*userGorm, error) {
	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &userGorm{db: db}, nil
}

// Close method closes the UserService database connection.
func (ug *userGorm) Close() error {
	return ug.db.Close()
}

// AutoMigrate method will attempt to automatically migrate the users table.
func (ug *userGorm) AutoMigrate() error {
	if err := ug.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}

// DestructiveReset method drops the user table and rebuilds it. Used in development env.
func (ug *userGorm) DestructiveReset() error {
	err := ug.db.DropTableIfExists(&User{}).Error
	if err != nil {
		return err
	}
	return ug.AutoMigrate()
}

func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
}

func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

func (ug *userGorm) Delete(id uint) error {
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}

// As a general rule, any error but ErrNotFound should probably result in a 500 error.
func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ByRemember looks up a user with the given remember token (which must be already hashed)
// and returns that user.
func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var user User
	err := first(ug.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}
