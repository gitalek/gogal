package models

import (
	"github.com/gitalek/gogal/hash"
	"github.com/gitalek/gogal/rand"
	"golang.org/x/crypto/bcrypt"
)

// userValidator is a validation layer that validates and normalizes
// data before passing it on the next UserDB in our interface chain.
type userValidator struct {
	UserDB
	hmac hash.HMAC
}

func (uv *userValidator) ByRemember(token string) (*User, error) {
	rememberHash := uv.hmac.Hash(token)
	return uv.UserDB.ByRemember(rememberHash)
}

func (uv *userValidator) Create(user *User) error {
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""

	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)

	return uv.UserDB.Create(user)
}

func (uv *userValidator) Update(user *User) error {
	if user.Remember != "" {
		user.RememberHash = uv.hmac.Hash(user.Remember)
	}
	return uv.Update(user)
}

func (uv *userValidator) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	return uv.Delete(id)
}
