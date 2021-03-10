package models

// userValidator is a validation layer that validates and normalizes
// data before passing it on the next UserDB in our interface chain.
type userValidator struct {
	UserDB
}

