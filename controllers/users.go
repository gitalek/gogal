package controllers

import (
	"fmt"
	"github.com/gitalek/gogal/views"
	"net/http"
)

type Users struct {
	NewView *views.View
}

func NewUsers() (*Users, error) {
	v, err := views.NewView("bootstrap", "users/new")
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
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "Email is %s\nPassword is %s\n", form.Email, form.Password)
}

type SignupForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}
