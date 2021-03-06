package controllers

import (
	"github.com/gitalek/gogal/views"
	"net/http"
)

type Users struct {
	NewView *views.View
}

func NewUsers() (*Users, error) {
	v, err := views.NewView("bootstrap", "views/users/new.gohtml")
	if err != nil {
		return nil, err
	}
	return &Users{NewView: v}, nil
}

// New processes the GET /new route
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

// Create processes the POST /signup route
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is a temporary response."))
}
