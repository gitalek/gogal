package models

import "github.com/gitalek/gogal/hash"

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
